package server

import (
	"fmt"
	"net/http"
	lobbystore "two-towers/lib/lobbyStore"

	"github.com/gorilla/pat"
)

func Run() error {

	store, err := NewAuth()
	if err != nil {
		return err
	}
	ls := lobbystore.NewLobbyStore()
	server := pat.New()

	RegisterRoutes(server, &ls, store)

	fmt.Println("running server")
	return http.ListenAndServe("localhost:8000", server)

}
