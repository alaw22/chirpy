package main

import (
	"io"
	"time"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request){
	defer req.Body.Close()

	type userEmail struct{
		Email string `json:"email"`
	}

	type userAccountInfo struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}

	// Read in request body
	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(w,
						  501,
						  "Couldn't read request body",
						  err)
		return
	}

	// get email from request
	email := userEmail{}

	err = json.Unmarshal(data,&email)
	if err != nil {
		respondWithError(w,
						  502,
						  "Couldn't unmarshal new user data",
						  err)
		return 
	}

	// Create user
	newUser, err := cfg.db.CreateUser(req.Context(),email.Email)
	if err != nil {
		respondWithError(w,
						  503,
						  "Error in CreateUser(), unable to create new user",
						  err)
		return
	}

	// Transfer new user data to encodable struct
	newUserInfo := userAccountInfo{
		ID: newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email: newUser.Email,
	}

	// Marshal new data to bytes
	respondWithJSON(w,
					201,
					newUserInfo)	


}