package main

import (
	"net/http"

	"github.com/alaw22/chirpy/internal/auth"
)


func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, req *http.Request) {

	type accessTokenString struct{
		Token string `json:"token"`
	}

	refreshTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w,
						 400,
						 "Unable to get token: Bearer token must not be formatted correctly",
						 err)
		return
	}

	// Lookup refresh token
	userID, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w,
						 401,
						 "Refresh token either revoked or doesn't exist",
						 err)

		return
	}

	if !userID.Valid {
		respondWithError(w,
						 500,
						 "User id NULL in refresh token table",
						 err)

		return		
	}

	// Create valid access token and send it to user
	tokenString, err := auth.MakeJWT(userID.UUID, cfg.secretKey, defaultExpirationTime)
	if err != nil {
		respondWithError(w,
						 500,
						 "Unable to create new refresh token",
						 err)
		return
	}

	// Send access token back to user
	respondWithJSON(w, 200, accessTokenString{
		Token: tokenString,
	})

}