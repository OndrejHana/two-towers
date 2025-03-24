package server

import (
	"fmt"
	"net/http"
	lobbystore "two-towers/lib/lobbyStore"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/markbates/goth/gothic"
)

var upgrader = websocket.Upgrader{} // use default options

func getAuthLogoutFunction(res http.ResponseWriter, req *http.Request) {
	gothic.Logout(res, req)
	http.Redirect(res, req, "/", http.StatusPermanentRedirect)
}

func getAuthProviderFunction(res http.ResponseWriter, req *http.Request) {
	if _, err := gothic.CompleteUserAuth(res, req); err != nil {
		gothic.BeginAuthHandler(res, req)
	}
}

func RegisterRoutes(router *pat.Router, ls *lobbystore.LobbyStore, sessionStore *sessions.CookieStore) {
	router.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Println(w, r, err)
		}

		fmt.Println(StoreUserSession(w, r, user))

		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	})
	router.Get("/auth/{provider}/logout", getAuthLogoutFunction)
	router.Get("/auth/{provider}", getAuthProviderFunction)
	router.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		user, err := GetSessionUser(r)
		fmt.Println(user, err)
	})

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("dist"))).Methods("GET")
	http.Handle("/", router)
}
