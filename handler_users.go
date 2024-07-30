package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/codezera11/chirps/internal/auth"
)

type User struct {
	Id int`json:"id"`
	Email string`json:"email"`
	IsChirpyRed bool`json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email string`json:"email"`
		Password string`json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := reqBody{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding params")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		Email: user.Email,
		Id: user.ID,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email string`json:"email"`
		Password string`json:"password"`
		ExpiresInSeconds int`json:"expires_in_seconds"`
	}

	type Response struct {
		User
		Token string`json:"token"`
		RefreshToken string`json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqBody{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding params")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error finding user")
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong password!")
	}

	defaultExpirationTime := 60 * 60
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpirationTime
	} else if params.ExpiresInSeconds > defaultExpirationTime {
		params.ExpiresInSeconds = defaultExpirationTime
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(defaultExpirationTime) * time.Second)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating jwt")
		return
	}

	refreshToken, err := auth.MakeRefToken();

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return 
	}

	refTokenExpTime := time.Now().UTC().Add(60 * 60 * 60)

	// add refresh token to user
	err = cfg.DB.AddRefTokenToUser(user.ID, refreshToken, refTokenExpTime)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		User: User{
			Id: user.ID,
			Email: user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token: jwt,
		RefreshToken: refreshToken,
	})
}