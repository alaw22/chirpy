package auth

import (
	"net/http"
	// "fmt"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error){
	authorization := headers.Get("Authorization")
	// if authorization == "" {
	// 	return "", fmt.Errorf("No authorization token")
	// }

	tokenString := strings.TrimSpace(strings.TrimPrefix(authorization,"Bearer"))
	
	return tokenString, nil
}