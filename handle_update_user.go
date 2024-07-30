package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/codezera11/chirps/internal/auth"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email string`json:"email"`
		Password string`json:"password"`
	}

	type Response struct {
		User
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	idString, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil  {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	reqParams := Request{}
	err = decoder.Decode(&reqParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(reqParams.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.DB.UpdateUser(reqParams.Email, hashedPassword, id)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		User: User{
			Email: user.Email,
			Id: user.ID,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}