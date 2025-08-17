package main

import (
	"io"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"time"

	"github.com/alaw22/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()


	type userCredentials struct{
		Password string `json:"password"`
		Email string `json:"email"`
	}

	type userAccountInfo struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(w,
		                 515,
						 "Error in io.ReadAll",
						 err)
		return
	}

	userCreds := userCredentials{}

	err = json.Unmarshal(data, &userCreds)
	if err != nil {
		respondWithError(w,
						 511,
						 "Error in loginHandler: Unable to unmarshal request body",
						 err)
		return
	}

	user, err := cfg.db.GetUserFromEmail(req.Context(),userCreds.Email)
	if err != nil {
		respondWithError(w,
					     401,
						 "Incorrect email or password",
						 err)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword,userCreds.Password)
	if err != nil {
		respondWithError(w,
		                 401,
						 "Incorrect email or password",
						 err)
		return
	}

	respondWithJSON(w,200,userAccountInfo{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})

}