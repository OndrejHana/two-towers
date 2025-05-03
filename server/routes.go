package server

import (
	"net/http"

	"github.com/OndrejHana/two-towers/server/handlers"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type CreateLobbyResponse struct {
	Code string `json:"code"`
}

type JoinLobbyRequest struct {
	Code     string `json:"code"`
	Username string `json:"username"`
}

type CreateLobbyRequest struct {
	Username string `json:"username"`
}

func SetupRoutes(router *pat.Router) {
	router.Handle("/lobby/new", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(handlers.NewLobbyHandler)))
	router.Handle("/lobby/join", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(handlers.JoinLobbyHandler)))
	router.Handle("/lobby/{code}", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(handlers.GetLobbyHandler)))
	router.Handle("/lobby/{code}/status", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(handlers.LobbyStatusHandler)))
	router.Handle("/lobby/{code}/leave", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(handlers.LeaveLobbyHandler)))
	router.Handle("/lobby/{code}/toggle", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(handlers.ToggleReadyHandler)))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("dist"))).Methods("GET")
}
