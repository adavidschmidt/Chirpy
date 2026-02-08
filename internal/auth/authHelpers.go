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
		return "", errors.New("Improper token prefix")
	}
	token := strings.TrimSpace(tokenString[len(prefix):])
	if token == "" {
		return "", errors.New("Empty bearer token")
	}
	return token, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("No API key provided")
	}
	const prefix = "ApiKey "
	if !strings.HasPrefix(tokenString, prefix) {
		return "", errors.New("Improper key prefix")
	}
	key := strings.TrimSpace(tokenString[len(prefix):])
	if key == "" {
		return "", errors.New("Empty api key")
	}
	return key, nil
}
