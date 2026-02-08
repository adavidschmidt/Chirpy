package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/adavidschmidt/Chirpy/internal/auth"
	"github.com/adavidschmidt/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type userLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := userLogin{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error deccding user information: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Issues with the provided password")
		return
	}

	createdUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:          createdUser.ID,
			CreatedAt:   createdUser.CreatedAt,
			UpdatedAt:   createdUser.UpdatedAt,
			Email:       createdUser.Email,
			IsChirpyRed: createdUser.IsChirpyRed,
		},
	})
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type userLogin struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		ExpiresIn *int   `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := userLogin{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad login data")
		return
	}
	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Cannot find user")
		return
	}
	goodPassword, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error checking password")
		return
	}
	if !goodPassword {
		respondWithError(w, http.StatusUnauthorized, "Password is incorrect")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting refresh token")
		return
	}

	tokenString, err := cfg.db.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: tokenString.Token,
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User User
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not provided")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid access token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bad request body")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		ID:             userID,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user data")
		return
	}

	returnUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, returnUser)
}
