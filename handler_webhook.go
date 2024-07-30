package main

import (
	"encoding/json"
	"net/http"

	"github.com/codezera11/chirps/internal/auth"
)

func (cfg *apiConfig) handlerWebook(w http.ResponseWriter, r *http.Request) {

	_, err := auth.GetApiKey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Api key not found")
		return
	}

	type request struct {
		Event string`json:"event"`
		Data struct{
			UserId int `json:"user_id"`
		}`json:"data"`
	}

	decoder := json.NewDecoder(r.Body)

	params := request{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding params")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, http.Response{})
		return
	}

	_, err = cfg.DB.GetUserById(params.Data.UserId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found!")
		return
	}

	err = cfg.DB.UpdateMembership(params.Data.UserId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating membership!")
		return
	}

	respondWithJSON(w, http.StatusNoContent, http.Response{})
}