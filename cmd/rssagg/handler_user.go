package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/natac13/bootdev-rssagg/internal/database"
)

func (a *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	var user database.CreateUserParams
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user.ID = uuid.New()
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	createdUser, err := a.DB.CreateUser(r.Context(), user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseUserToAPIUser(createdUser))
}

func (a *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToAPIUser(user))
}
