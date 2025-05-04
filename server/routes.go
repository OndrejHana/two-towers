package server

import (
	"net/http"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SetupRoutes(router *pat.Router) {
	router.Handle("/lobby/new", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(NewLobbyHandler)))
	router.Handle("/lobby/join", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(JoinLobbyHandler)))
	router.Handle("/lobby/{code}", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(GetLobbyHandler)))
	router.Handle("/lobby/{code}/status", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(LobbyStatusHandler)))
	router.Handle("/lobby/{code}/leave", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(LeaveLobbyHandler)))
	router.Handle("/lobby/{code}/toggle", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(ToggleReadyHandler)))
	router.Handle("/game/{gameId}/ws", http.HandlerFunc(HandleGameWs))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("dist"))).Methods("GET")
}
