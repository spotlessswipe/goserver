package main

import (
	"goserver/internal/auth"
	"goserver/internal/database"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	path := r.PathValue("chirpID")
	id, err := uuid.Parse(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing path to UUID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Forbidden", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     id,
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
