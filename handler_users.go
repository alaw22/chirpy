package main

import (
	"io"
	"time"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"

	"github.com/alaw22/chirpy/internal/auth"
	"github.com/alaw22/chirpy/internal/database"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request){
	defer req.Body.Close()

	type userCredentials struct{
		Email string `json:"email"`
		Password string `json:"password"`
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
	creds := userCredentials{}

	err = json.Unmarshal(data,&creds)
	if err != nil {
		respondWithError(w,
						 502,
						 "Couldn't unmarshal new user data",
						 err)
		return 
	}

	hashedPassword, err := auth.HashPassword(creds.Password)
	if err != nil {
		respondWithError(w,
						 510,
						 "Couldn't hash password",
						 err)
		return
	}

	// Create user
	newUser, err := cfg.db.CreateUser(req.Context(),
									  database.CreateUserParams{
										Email: creds.Email,
										HashedPassword: hashedPassword,
									  })

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