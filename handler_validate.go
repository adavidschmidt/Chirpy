package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"unicode/utf8"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if utf8.RuneCountInString(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := helperCleanBody(params.Body)

	respondWithJSON(w, http.StatusOK, returnVal{
		CleanedBody: cleanedBody,
	})
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
