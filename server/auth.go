package server

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

const (
	key    = "randomString"
	MaxAge = 86400
	IsProd = false
)

func NewAuth() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	GHClientID := os.Getenv("GITHUB_ID")
	GHClientSecret := os.Getenv("GITHUB_SECRET")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store

	goth.UseProviders(github.New(GHClientID, GHClientSecret, "http://localhost:8000/auth/github/callback"))
}
