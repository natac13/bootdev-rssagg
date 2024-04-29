package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/natac13/bootdev-rssagg/internal/database"
)

func (a *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var feed struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&feed)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
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
	follow, err := a.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    createdFeed.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = follow
	respondWithJSON(w, http.StatusCreated, databaseFeedToAPIFeed(createdFeed))
}

func (a *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := a.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	apiFeeds := make([]Feed, len(feeds))
	for i, feed := range feeds {
		apiFeeds[i] = databaseFeedToAPIFeed(feed)
	}
	respondWithJSON(w, http.StatusOK, apiFeeds)

}
