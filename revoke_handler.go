package main

import (
	"net/http"

	"github.com/alaw22/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeRefreshHandler(w http.ResponseWriter, req *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w,
						 400,
						 "Unable to get token: Bearer token must not be formatted correctly",
						 err)
		return
	}

	err = cfg.db.RevokeToken(req.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w,
						 500,
						 "Token doesn't exist in database",
						err)
		return
	}

	// send status code
	w.WriteHeader(204)

}