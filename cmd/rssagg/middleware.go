package main

import (
	"net/http"

	"github.com/natac13/bootdev-rssagg/internal/auth"
	"github.com/natac13/bootdev-rssagg/internal/database"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type authenticatedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (cfg *apiConfig) authMiddleware(handler authenticatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey, err := auth.GetApiKey(r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		handler(w, r, user)
	}
}
