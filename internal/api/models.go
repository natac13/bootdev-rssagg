package api

import "github.com/natac13/bootdev-rssagg/internal/database"

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
