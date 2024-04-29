package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/natac13/bootdev-rssagg/internal/auth"
	"github.com/natac13/bootdev-rssagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func NewAPI(db *database.Queries) *apiConfig {
	return &apiConfig{
		DB: db,
	}
}

func (a *apiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {

	var user database.CreateUserParams
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user.ID = uuid.New()
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	createdUser, err := a.DB.CreateUser(r.Context(), user)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, databaseUserToAPIUser(createdUser))
}

func (a *apiConfig) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := a.DB.GetUserByApiKey(r.Context(), apiKey)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, databaseUserToAPIUser(user))
}
