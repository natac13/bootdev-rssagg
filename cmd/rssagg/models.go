package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/natac13/bootdev-rssagg/internal/database"
)

type User struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Apikey    string `json:"apikey"`
}

func databaseUserToAPIUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID.String(),
		CreatedAt: dbUser.CreatedAt.String(),
		UpdatedAt: dbUser.UpdatedAt.String(),
		Name:      dbUser.Name,
		Apikey:    dbUser.Apikey,
	}
}

type Feed struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	UserID    string `json:"user_id"`
}

func databaseFeedToAPIFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID.String(),
		CreatedAt: dbFeed.CreatedAt.String(),
		UpdatedAt: dbFeed.UpdatedAt.String(),
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		UserID:    dbFeed.UserID.String(),
	}
}

func databaseFeedsToAPIFeeds(dbFeeds []database.Feed) []Feed {
	apiFeeds := make([]Feed, len(dbFeeds))
	for i, dbFeed := range dbFeeds {
		apiFeeds[i] = databaseFeedToAPIFeed(dbFeed)
	}
	return apiFeeds
}

type FeedFollow struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UserID    string `json:"user_id"`
	FeedID    string `json:"feed_id"`
}

func databaseFeedFollowToAPIFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID.String(),
		CreatedAt: dbFeedFollow.CreatedAt.String(),
		UpdatedAt: dbFeedFollow.UpdatedAt.String(),
		UserID:    dbFeedFollow.UserID.String(),
		FeedID:    dbFeedFollow.FeedID.String(),
	}
}

func databaseFeedFollowsToAPIFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	apiFeedFollows := make([]FeedFollow, len(dbFeedFollows))
	for i, dbFeedFollow := range dbFeedFollows {
		apiFeedFollows[i] = databaseFeedFollowToAPIFeedFollow(dbFeedFollow)
	}
	return apiFeedFollows
}

type Post struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	FeedID      uuid.UUID `json:"feed_id"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
}

func databasePostToAPIPost(dbPost database.Post) Post {
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}
	return Post{
		ID:          dbPost.ID.String(),
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		FeedID:      dbPost.FeedID,
		Title:       dbPost.Title,
		Url:         dbPost.Url,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
	}
}

func databasePostsToAPIPosts(dbPosts []database.Post) []Post {
	apiPosts := make([]Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		apiPosts[i] = databasePostToAPIPost(dbPost)
	}
	return apiPosts
}
