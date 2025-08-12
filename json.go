package main

import (
	"net/http"
	"encoding/json"
	"log"
)

func respondeWithError(w http.ResponseWriter, statusCode int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if statusCode > 499 {
		log.Printf("Responding with 5XX error: %s\n",msg)
	}

	type jsonError struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, statusCode, jsonError{
		Error: msg,
	})
}


func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type","application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Couldn't marshal payload: %w\n",err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(data)
	
}