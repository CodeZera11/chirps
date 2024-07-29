package main

import (
	"net/http"

	"github.com/codezera11/chirps/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string`json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	jwtToken, err := cfg.DB.GetNewAccessToken(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, response{
			Token: jwtToken,
	})
}