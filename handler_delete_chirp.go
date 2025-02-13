package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codezera11/chirps/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found")
		return
	}

	idString, err := auth.ValidateJWT(token, cfg.jwtSecret)

	fmt.Println("id string here", idString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error validating token")
		return
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error converting id")
		return
	}

	chirpIdStr := r.PathValue("chirpId")

	chirpId, err := strconv.Atoi(chirpIdStr)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Please enter a valid chirp id!")
		return
	}

	dbChirp, err := cfg.DB.GetOneChirp(chirpId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error finding chirp")
		return
	}

	if dbChirp.AuthorId != id {
		respondWithError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	_, err = cfg.DB.DeleteChirp(chirpId, id)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp")
		return
	}

	respondWithJSON(w, http.StatusNoContent, http.Response{})
}