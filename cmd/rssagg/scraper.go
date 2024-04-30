package main

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/natac13/bootdev-rssagg/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	slog.Info("Scraping on feeds started", "concurrency", concurrency, "timeBetweenRequest", timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			slog.Error("Failed to get feeds to fetch", "error", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		slog.Error("Failed to mark feed as fetched", "error", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		slog.Error("Failed to fetch feed", "error", err)
		return
	}

	slog.Info("scarping feed", "feed", feed.Url, "count", len(rssFeed.Channel.Item))

	for _, item := range rssFeed.Channel.Item {
		des := sql.NullString{}
		if item.Description != "" {
			des.String = item.Description
			des.Valid = true
		}
		pub_at, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			slog.Error("Failed to parse time", "error", err, "time", item.PubDate)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			FeedID:      feed.ID,
			Title:       item.Title,
			Url:         item.Link,
			Description: des,
			PublishedAt: pub_at,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				continue
			}
			slog.Error("Failed to create post", "error", err)
			continue
		}
	}

	slog.Info("Feed fetched", "feed", feed.Url, "count", len(rssFeed.Channel.Item))
}
