package main

import (
	"io"
	"encoding/json"
	"net/http"
)

func validateChirpHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	
	type chirpBody struct {
		Body string `json:"body"`
	}
	
	type jsonValid struct {
		Valid bool `json:"valid"`
	}
	
	dat, err := io.ReadAll(req.Body)
	if err != nil {
		respondeWithError(w,
						  501,
						  "Couldn't read request body",
						  err)
		return
	}


	chirp := chirpBody{}


	err = json.Unmarshal(dat,&chirp)
	if err != nil {

		respondeWithError(w,
						  502,
						  "Couldn't unpack json to chirpBody{}",
						  err)
						  
		return
	}

	if len(chirp.Body) > 140 {
		respondeWithError(w,
						  400,
						  "Chirp is too long",
						  nil)
	} else {

		respondWithJSON(w, 200, jsonValid{
			Valid: true,
		})
		
	}

}