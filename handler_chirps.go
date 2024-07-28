package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Chirp struct {
	ID int`json:"id"`
	Body string`json:"body"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	type reqBody struct {
		Body string`json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqBody{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 500, "Error decoding request body")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedInput, err := validateInput(params.Body)

	if err != nil {
		respondWithError(w, 500, "Error validating body")
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleanedInput)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		Body: chirp.Body,
	})
}

func validateInput(val string) (string, error) {

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	if len(val) == 0 {
		return "", errors.New("please send a message") 
	}

	words := strings.Split(val, " ")

	for i, word := range words {
		for _, pWord := range profaneWords {
			lWord := strings.ToLower(word)
			if lWord == pWord {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " "), nil
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()


	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return 
	}

	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID: dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetOneChirp(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Please enter a valid id!")
		return
	}

	chirp, err := cfg.DB.GetOneChirp(id)

	if errors.Is(err, os.ErrNotExist) {
		respondWithError(w, http.StatusNotFound, "Not found")
	}

	respondWithJSON(w, http.StatusOK, chirp)
}