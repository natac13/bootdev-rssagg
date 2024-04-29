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

func (a *apiConfig) HandleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	RespondWithJSON(w, http.StatusOK, databaseUserToAPIUser(user))
}

func (a *apiConfig) HandleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var feed struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&feed)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	createdFeed, err := a.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		Url:       feed.URL,
		Name:      feed.Name,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ID:        uuid.New(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, databaseFeedToAPIFeed(createdFeed))
}

func (a *apiConfig) HandleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := a.DB.GetFeeds(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	apiFeeds := make([]Feed, len(feeds))
	for i, feed := range feeds {
		apiFeeds[i] = databaseFeedToAPIFeed(feed)
	}
	RespondWithJSON(w, http.StatusOK, apiFeeds)

}

// MIDDLEWARE
type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (cfg *apiConfig) AuthMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey, err := auth.GetApiKey(r)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		handler(w, r, user)
	}
}
