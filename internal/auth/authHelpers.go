package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("No auth token provided")
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(tokenString, prefix) {
		return "", errors.New("Improper toekn prefix")
	}
	token := strings.TrimSpace(tokenString[len(prefix):])
	if token == "" {
		return "", errors.New("Empty bearer token")
	}
	return token, nil
}
