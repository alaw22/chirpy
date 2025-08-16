package main

import (
	"io"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"github.com/google/uuid"
	"github.com/alaw22/chirpy/internal/database"
)

func replace_profanity(msg string) string{
	profane_words := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	
	split_msg := strings.Split(msg," ")

	for i, word := range split_msg{
		loweredWord := strings.ToLower(word)

		if _, ok := profane_words[loweredWord]; ok{
			split_msg[i] = "****"
		}
	}

	return strings.Join(split_msg," ")

}



func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	
	type chirpInfoFull struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	
	type chirpInfoStripped struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	// type cleanChirp struct {
	// 	Clean_Body string `json:"cleaned_body"`
	// }
	

	dat, err := io.ReadAll(req.Body)
	if err != nil {
		respondeWithError(w,
						  501,
						  "Couldn't read request body",
						  err)
		return
	}


	chirp := chirpInfoStripped{}


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

		// Create chirp params struct with cleaned chirp body
		chirpParams := database.CreateChirpParams{
			Body: replace_profanity(chirp.Body),
			UserID: chirp.UserID,
		}

		// Store chirp before responding
		chirpEntry, err := cfg.db.CreateChirp(req.Context(), chirpParams)
		if err != nil {
			respondeWithError(w,
							  506,
							  "Couldn't create chirp",
							  err)
			return
		}

		// Marshal new data
		respondWithJSON(w, 201, chirpInfoFull{
			ID: chirpEntry.ID,
			CreatedAt: chirpEntry.CreatedAt,
			UpdatedAt: chirpEntry.UpdatedAt,
			Body: chirpEntry.Body,
			UserID: chirpEntry.UserID,
		})
		
	}

}