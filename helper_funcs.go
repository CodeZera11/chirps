package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {

	type ErrorType struct {
		Error string`json:"error"`
	}

	w.WriteHeader(code)
	error := ErrorType{
		Error: msg,
	}

	data, err := json.Marshal(error)

	if err != nil {
		fmt.Printf("Error marshalling: %s\n", err)
		return
	}

	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)

	data, err := json.Marshal(payload)

	if err != nil {
		fmt.Printf("Error marshalling: %s\n", err)
		return
	}

	w.Write(data)
}