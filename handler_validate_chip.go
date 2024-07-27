package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type respType struct {
		CleanedBody string `json:"cleaned_body"`
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedInput, err := validateInput(params.Body)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	resp := respType{
		CleanedBody: cleanedInput,
	}

	respondWithJSON(w, 200, resp)
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