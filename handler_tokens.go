package main

import (
	"net/http"
	"time"

	"github.com/adavidschmidt/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "refresh token not provided")
		return
	}
	token, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusBadGateway, "refresh token not found")
		return
	}
	if token.ExpiresAt.Before(time.Now().UTC()) || token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token not found or expired")
		return
	}

	refreshedToken, err := auth.MakeJWT(token.UserID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making new token")
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: refreshedToken,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "token not provided")
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error revoking token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
