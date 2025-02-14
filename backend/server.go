package backend

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"two-towers/backend/router"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func setupAuth() error {
	key := os.Getenv("SESSION_STORE_KEY")

	maxAge, err := strconv.Atoi(os.Getenv("SESSION_MAX_AGE"))
	if err != nil {
		return err
	}

	isProd, err := strconv.ParseBool(os.Getenv("IS_PROD"))
	if err != nil {
		return err
	}

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:8080/auth/callback"),
	)

	return nil
}

func Serve() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	if err := setupAuth(); err != nil {
		return err
	}

	mux := http.NewServeMux()
	backend.RegisterRoutes(mux)

	port := os.Getenv("PORT")
	return http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
