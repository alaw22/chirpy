package auth

import (
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error){
	authorization := headers.Get("Authorization")

	tokenString := strings.TrimSpace(strings.TrimPrefix(authorization,"Bearer"))
	
	return tokenString, nil
}