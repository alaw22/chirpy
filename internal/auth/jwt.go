package auth

import (
	// "fmt"
	"time"

	uuid "github.com/google/uuid"
	jwt "github.com/golang-jwt/jwt/v5"
)


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Create claims
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	signedJWT, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Validate tokenString
	token, err := jwt.ParseWithClaims(tokenString,
									  &jwt.RegisteredClaims{},
									  func(token *jwt.Token) (any, error){
										return []byte(tokenSecret), nil
									})

	if err != nil {
		return uuid.UUID{}, err
	}

	// Get user id as a string from token claims
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	// Transform userIDString into a uuid.UUID instance
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.UUID{}, err
	}


	return userID, nil

} 

