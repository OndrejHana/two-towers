package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"two-towers/lib/lobbyStore"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

type User struct {
	ID string
}

// func getAuth(r *http.Request) (*clerk.User, bool, error) {
// 	ctx := r.Context()
// 	claims, ok := clerk.SessionClaimsFromContext(ctx)
// 	if !ok {
// 		return nil, false, nil
// 	}
//
// 	u, err := user.Get(ctx, claims.Subject)
// 	if err != nil {
// 		return nil, ok, err
// 	}
//
// 	return u, true, nil
// }

func getAuth(r *http.Request) (*User, bool, error) {
	return nil, false, errors.New("Not implemented")
}

func RegisterRoutes(mux *http.ServeMux, ls *lobbystore.LobbyStore) {
	mux.Handle("/", http.FileServer(http.Dir("dist")))

	newLobbyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok, err := getAuth(r)
		fmt.Println("auth: ", user, ok, err)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		l := ls.NewLobby(user.ID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lobbystore.NewLobbyRes{
			ConnString: l.GetConnString(),
		})
		fmt.Println("done")
	})

	mux.Handle("/api/lobby/new", newLobbyHandler)

	// lobbyPlayersHandler := http.HandlerFunc()

	mux.HandleFunc("/api/lobby/players", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("got here")

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ws.Close()

		if err := ws.WriteMessage(websocket.TextMessage, []byte("sup")); err != nil {
			fmt.Println("sup", err)
		}
		if err := ws.WriteMessage(websocket.TextMessage, []byte("bob")); err != nil {
			fmt.Println("sup", err)
		}
		if err := ws.WriteMessage(websocket.TextMessage, []byte("jon")); err != nil {
			fmt.Println("sup", err)
		}
	})
}
