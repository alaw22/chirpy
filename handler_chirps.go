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

type chirpInfoFull struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

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
	
	type chirpInfoStripped struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	dat, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(w,
						  501,
						  "Couldn't read request body",
						  err)
		return
	}


	chirp := chirpInfoStripped{}


	err = json.Unmarshal(dat,&chirp)
	if err != nil {

		respondWithError(w,
						  502,
						  "Couldn't unpack json to chirpBody{}",
						  err)
						  
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w,
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
			respondWithError(w,
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

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w,
						  507,
						  "Error in GetAllChirps(): Couldn't get all chirps",
						  err)
		return
	}

	// Unpack into encodable slice of structs 
	chirpsInfo := make([]chirpInfoFull,len(chirps))
	for i, chirp := range chirps{
		chirpsInfo[i].ID = chirp.ID
		chirpsInfo[i].CreatedAt = chirp.CreatedAt
		chirpsInfo[i].UpdatedAt = chirp.UpdatedAt
		chirpsInfo[i].Body = chirp.Body
		chirpsInfo[i].UserID = chirp.UserID
	} 

	respondWithJSON(w,200,chirpsInfo)
}