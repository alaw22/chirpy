package main

import (
	"io"
	"fmt"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"sort"
	"github.com/google/uuid"
	
	"github.com/alaw22/chirpy/internal/database"
	"github.com/alaw22/chirpy/internal/auth"
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
	}

	// Authenticate user
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w,
					     401,
						 "Unauthorized",
						 err)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.secretKey)
	if err != nil {
		respondWithError(w,
						 401,
						 "Unauthorized user",
						 err)
		return
	}


	dat, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(w,
						  500,
						  "Couldn't read request body",
						  err)
		return
	}


	chirp := chirpInfoStripped{}


	err = json.Unmarshal(dat,&chirp)
	if err != nil {

		respondWithError(w,
						 500,
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
			UserID: userID,
		}

		// Store chirp before responding
		chirpEntry, err := cfg.db.CreateChirp(req.Context(), chirpParams)
		if err != nil {
			respondWithError(w,
							500,
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
	var chirps []database.Chirp
	var err error
	var authorID uuid.UUID

	// check to see if an author_id was passed to get all chirps of a author author_id
	authorIDString := req.URL.Query().Get("author_id")
	if authorIDString != ""{

		authorID, err = uuid.Parse(authorIDString)
		if err != nil{
			respondWithError(w,
							 400,
							 "Unable to create uuid.UUID from author ID provided",
							 err)
			return
		}

		chirps, err = cfg.db.GetAllChirpsByUser(req.Context(), authorID)
		if err != nil {
			respondWithError(w,
							 500,
							 "Error in GetAllChirpsByUser(): Couldn't get chirps",
							 err)
			return
		}
		
	} else {
		
		chirps, err = cfg.db.GetAllChirps(req.Context())
		if err != nil {
			respondWithError(w,
							 500,
							 "Error in GetAllChirps(): Couldn't get all chirps",
							 err)
			return
		}
	}

	sortQ := req.URL.Query().Get("sort")
	if sortQ == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt)})
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

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, req *http.Request) {
	chirpIDString := req.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w,
						 500,
						 fmt.Sprintf("%s isn't a UUID something is really wrong",chirpIDString),
						 err)
		return
	}

	// get chirp with id == chirpID
	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w,
						 404,
						 "Error in GetChirp(): unable to get chirp",
						 err)
		return
	}

	respondWithJSON(w, 200, chirpInfoFull{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, req *http.Request) {
	accessTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil{
		respondWithError(w,
						 400,
						 "Unable to get access token",
						 err)
		return
	}
	
	
	userID, err := auth.ValidateJWT(accessTokenString, cfg.secretKey)
	if err != nil{
		respondWithError(w,
						 401,
						 "Not a valid user",
						 err)
		return
	}

	chirpIDString := req.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w,
						 500,
						 fmt.Sprintf("%s isn't a UUID something is really wrong",chirpIDString),
						 err)
		return
	}
	
	// get chirp with id == chirpID
	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w,
						 404,
						 "Error in GetChirp(): unable to get chirp",
						 err)
		return
	}

	// check author is the same as the authenticated user. Otherwise known as 
	// authorization
	if chirp.UserID != userID {
			respondWithError(w,
							403,
							"You are not the author of the chirp you are trying to delete",
							err)
		return		
	}

	// Delete chirp
	err = cfg.db.DeleteChirp(req.Context(), chirp.ID)
	if err != nil{
		respondWithError(w,
						 500,
						 "Unable to delete chirp",
						 err)
		return
	}

	w.WriteHeader(204)

}