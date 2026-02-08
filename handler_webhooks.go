package main

import (
	"encoding/json"
	"net/http"

	"github.com/adavidschmidt/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}

	if apiKey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := payload{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Prvoided user id is invalid")
		return
	}

	err = cfg.db.AddChirpyRed(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cannot find user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
