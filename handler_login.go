package main

import (
	"io"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"time"

	"github.com/alaw22/chirpy/internal/auth"
	"github.com/alaw22/chirpy/internal/database"
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
		AccessTokenString string `json:"token"`
		RefreshTokenString string `json:"refresh_token"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(w,
		                 500,
						 "Error in io.ReadAll",
						 err)
		return
	}

	userCreds := userCredentials{}

	err = json.Unmarshal(data, &userCreds)
	if err != nil {
		respondWithError(w,
						 500,
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

	// Create access token
	tokenString, err := auth.MakeJWT(user.ID, cfg.secretKey, defaultExpirationTime)
	if err != nil {
		respondWithError(w,
						 500,
						 "Unable to make JWT",
						 err)
		return
	}

	// Create refresh token
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w,
						 500,
						 "Unable to Create refresh token string",
						 err)
		return
	}

	// Store refresh token
	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
																  	Token: refreshTokenString,
																  	UserID: user.ID,
																  })

	if err != nil {
		respondWithError(w,
						 500,
						 "Unable to store refresh token",
						 err)
		return
	}

	respondWithJSON(w,200,userAccountInfo{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		AccessTokenString: tokenString,
		RefreshTokenString: refreshTokenString,
		IsChirpyRed: user.IsChirpyRed,
	})

}