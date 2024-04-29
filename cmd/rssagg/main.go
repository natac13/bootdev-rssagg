package main

import (
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	server := NewServer(":" + portStr)

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		panic("DATABASE_URL is required")
	}
	server.Start(dbUrl)

}
