package server

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func Init() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

	router := pat.New()
	SetupRoutes(router)

	log.Printf("Starting server on %s", *addr)
	return http.ListenAndServe(":8000", router)
}
