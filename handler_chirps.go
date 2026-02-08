package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/adavidschmidt/Chirpy/internal/auth"
	"github.com/adavidschmidt/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	Userid    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid bearer token")
		return
	}
	userId, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user")
	}
	decoder := json.NewDecoder(r.Body)
	params := chirpRequest{}
	err = decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if utf8.RuneCountInString(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := helperCleanBody(params.Body)
	chirpCreate, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userId,
	})

	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	chirp := Chirp{
		ID:        chirpCreate.ID,
		CreatedAt: chirpCreate.CreatedAt,
		UpdatedAt: chirpCreate.UpdatedAt,
		Body:      chirpCreate.Body,
		Userid:    chirpCreate.UserID,
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func helperCleanBody(s string) string {
	words := []string{"kerfuffle", "sharbert", "fornax"}
	sSplit := strings.Split(s, " ")
	for i, word := range sSplit {
		if slices.Contains(words, strings.ToLower(word)) {
			sSplit[i] = "****"
		}
	}
	return strings.Join(sSplit, " ")
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	authorID := uuid.Nil
	var err error
	if s == "" {
		authorID = uuid.Nil
	} else {
		authorID, err = uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author id")
		}
	}

	sortDir := r.URL.Query().Get("sort")
	data, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Problem getting chirps: %s", err)
		respondWithError(w, 500, "Error retreiving chirps")
	}
	returnChirps := []Chirp{}
	for _, chirp := range data {
		if chirp.UserID != authorID && authorID != uuid.Nil {
			continue
		}
		returnChirps = append(returnChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			Userid:    chirp.UserID,
		})

	}
	if sortDir != "" {
		if sortDir == "desc" {
			sort.Slice(returnChirps, func(i, j int) bool {
				return returnChirps[i].CreatedAt.After(returnChirps[j].CreatedAt)
			})
		} else {
			sort.Slice(returnChirps, func(i, j int) bool {
				return returnChirps[i].CreatedAt.Before(returnChirps[j].CreatedAt)
			})
		}
	}

	respondWithJSON(w, 200, returnChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Error converting provided id to uuid: %s", err)
		respondWithError(w, http.StatusBadRequest, "Bad chirp id provided")
	}
	chirp, err := cfg.db.GetChirp(r.Context(), uuid)
	if err != nil {
		log.Printf("Error getting chirp %s: %s", id, err)
		respondWithError(w, 404, "Unable to find chirp")
	}
	returnChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		Userid:    chirp.UserID,
	}
	respondWithJSON(w, 200, returnChirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Error converting provided id to uuid: %s", err)
		respondWithError(w, http.StatusBadRequest, "Bad chirp id provided")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "auth token not provided")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "invalid token")
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "did not find chirp")
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "user did not write chirp")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "did not find chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
