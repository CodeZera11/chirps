package main

import (
	"net/http"

	"github.com/codezera11/chirps/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	type response struct {

	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	err = cfg.DB.RevokeRefreshToken(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent,	response{})
}