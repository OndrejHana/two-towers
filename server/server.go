package server

import (
	"fmt"
	"net/http"
	lobbystore "two-towers/lib/lobbyStore"

	"github.com/gorilla/pat"
)

func Run() error {
	// err := godotenv.Load()
	// if err != nil {
	// 	return fmt.Errorf("Error loading .env file")
	// }
	// key := os.Getenv("CLERK_SECRET_KEY")
	// if key == "" {
	// 	return fmt.Errorf("CLERK_SECRET_KEY not set")
	// }
	// clerk.SetKey(key)

	NewAuth()
	ls := lobbystore.NewLobbyStore()
	server := pat.New()

	RegisterRoutes(server, &ls)

	fmt.Println("running server")
	return http.ListenAndServe("localhost:8000", server)

}
