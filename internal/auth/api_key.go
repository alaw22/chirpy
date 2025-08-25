package auth

import (
	"net/http"
	"strings"
)

// Exactly like GetBearerToken so yeah
func GetAPIKey(headers http.Header) (string, error){
	authorization := headers.Get("Authorization")

	keyString := strings.TrimSpace(strings.TrimPrefix(authorization,"ApiKey"))

	return keyString, nil
}