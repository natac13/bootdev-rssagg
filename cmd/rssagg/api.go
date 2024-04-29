package main

import (
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
