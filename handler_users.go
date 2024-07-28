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
	}

	decoder := json.NewDecoder(r.Body)

	params := reqBody{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode params")
	}

	user, err := cfg.DB.CreateUser(params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
	}

	respondWithJSON(w, http.StatusCreated, user)
}