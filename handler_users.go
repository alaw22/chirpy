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
		IsChirpyRed bool `json:"is_chirpy_red"`
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
		IsChirpyRed: newUser.IsChirpyRed,
	}

	// Marshal new data to bytes
	respondWithJSON(w,
					201,
					newUserInfo)	


}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	
	accessTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w,
						 401,
						 "Most likely a malformed token header",
						 err)
		return
	}

	// Validate user
	userID, err := auth.ValidateJWT(accessTokenString, cfg.secretKey)
	if err != nil{
		respondWithError(w,
						 401,
						 "Unable to validate jwt in user update",
						 err)
		return
	}

	type newUserCredentials struct{
		NewEmail string `json:"email"`
		NewPassword string `json:"password"`
	}

	type userAccountInfo struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(w,
						 500,
						 "Unable to read request body",
						 err)
		return
	}

	newUserCreds := newUserCredentials{}

	err = json.Unmarshal(data, &newUserCreds)
	if err != nil{
		respondWithError(w,
						 400,
						 "Unable to unmarshal new user credentials",
						 err)
		return
	}

	newHashedPassword, err := auth.HashPassword(newUserCreds.NewPassword)
	if err != nil{
		respondWithError(w,
						 500,
						 "Unable to hash new password",
						 err)
		return
	}

	updatedUserParams := database.UpdateUserEmailAndPasswordParams{
		Email: newUserCreds.NewEmail,
		HashedPassword: newHashedPassword,
		ID: userID,
	}

	// Update credentials
	updatedUser, err := cfg.db.UpdateUserEmailAndPassword(req.Context(),updatedUserParams)

	if err != nil{
		respondWithError(w,
						 500,
						 "Unable to update user credentials",
						 err)
		return
	}


	updatedUserInfo := userAccountInfo{
		ID: updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email: updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	}

	respondWithJSON(w,200,updatedUserInfo)


}


func (cfg *apiConfig) upgradeUserHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// Webhook struct
	type upgradeRequest struct {
		Event string `json:"event"`
		Data struct{
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	data, err := io.ReadAll(req.Body)
	if err != nil{
		respondWithError(w,
						 400,
						 "Unable to read request from webhook",
						 err)
		return
	}

	upgrade_req := upgradeRequest{}

	err = json.Unmarshal(data,&upgrade_req)
	if err != nil{
		respondWithError(w,
						 400,
						 "Ill formatted request in webhook",
						 err)
		return
	}

	if upgrade_req.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	userID, err := uuid.Parse(upgrade_req.Data.UserID)
	if err != nil{
		respondWithError(w,
						 400,
						 "Unable to convert provided user_id to UUID",
						 err)
		return
	}

	// upgrade user
	err = cfg.db.UpgradeUser(req.Context(), userID)
	if err != nil{
		respondWithError(w,
						 404,
						 "User cannot be found",
						 err)
		return
	}

	w.WriteHeader(204)
	

}