package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Id int`json:"id"`
	Email string`json:"email"`
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

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := cfg.DB.LoginUser(params.Email, params.Password)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, string(err.Error()))
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}