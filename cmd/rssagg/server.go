package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/natac13/bootdev-rssagg/internal/api"
	"github.com/natac13/bootdev-rssagg/internal/database"
)

type Server struct {
	server *http.Server
	// v1 router
	router     *http.ServeMux
	listenAddr string
}

func NewServer(listenAddr string) *Server {
	v1Router := http.NewServeMux()

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	server := http.Server{
		Addr:    listenAddr,
		Handler: middlewareCors(mainRouter),
	}

	return &Server{
		server:     &server,
		router:     v1Router,
		listenAddr: listenAddr,
	}
}

func (s *Server) Start(dbUrl string) {
	s.setupRoutes(dbUrl)

	fmt.Printf("Server is running on %s\n", s.listenAddr)
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func (s *Server) setupRoutes(dbUrl string) {
	fmt.Printf("Database URL: %s\n", dbUrl)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	apiConfig := api.NewAPI(dbQueries)

	s.router.HandleFunc("GET /readiness", handlerReadiness)
	s.router.HandleFunc("GET /error", func(w http.ResponseWriter, r *http.Request) {
		api.RespondWithError(w, http.StatusInternalServerError, "error")
	})

	s.router.HandleFunc("POST /users", apiConfig.HandleCreateUser)
	s.router.HandleFunc("GET /users", apiConfig.HandleGetUser)
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	api.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
