package server

import (
	"fmt"
	"net/http"
	lobbystore "two-towers/lib/lobbyStore"

	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
	"github.com/markbates/goth/gothic"
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

//	func getAuth(r *http.Request) (*User, bool, error) {
//		return nil, false, errors.New("Not implemented")
//	}

func getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Println(w, r, err)
	}

	fmt.Println("getAuthCallbackFunction", user)

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func getAuthLogoutFunction(res http.ResponseWriter, req *http.Request) {
	gothic.Logout(res, req)
	http.Redirect(res, req, "/", http.StatusPermanentRedirect)
	// res.Header().Set("Location", "/")
	// res.WriteHeader(http.StatusTemporaryRedirect)
}

func getAuthProviderFunction(res http.ResponseWriter, req *http.Request) {
	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		fmt.Println("getAuthProviderFunction", gothUser)
	} else {
		fmt.Println("authenticating")
		gothic.BeginAuthHandler(res, req)
	}
}

func RegisterRoutes(router *pat.Router, ls *lobbystore.LobbyStore) {
	router.Get("/auth/{provider}/callback", getAuthCallbackFunction)
	router.Get("/auth/{provider}/logout", getAuthLogoutFunction)
	router.Get("/auth/{provider}", getAuthProviderFunction)

	// router.Handle("/", http.FileServer(http.Dir("dist")))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("dist"))).Methods("GET")
	http.Handle("/", router)

	// mux.HandleFunc("/auth/callback", getAuthCallbackFunction)
	// mux.HandleFunc("/auth", getAuthProviderFunction)
	// mux.HandleFunc("/logout", getAuthLogoutFunction)

	// mux.Handle("/", http.FileServer(http.Dir("dist")))
	// mux.HandleFunc("/auth/callback", getAuthCallbackFunction)
	// mux.HandleFunc("/auth", getAuthProviderFunction)
	// mux.HandleFunc("/logout", getAuthLogoutFunction)
	// fmt.Println("registered")

	// fmt.Println("registered")

	// newLobbyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	user, ok, err := getAuth(r)
	// 	fmt.Println("auth: ", user, ok, err)
	//
	// 	if !ok {
	// 		w.WriteHeader(http.StatusUnauthorized)
	// 		return
	// 	}
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	l := ls.NewLobby(user.ID)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(lobbystore.NewLobbyRes{
	// 		ConnString: l.GetConnString(),
	// 	})
	// 	fmt.Println("done")
	// })
	//
	// mux.Handle("/api/lobby/new", newLobbyHandler)
	//
	// // lobbyPlayersHandler := http.HandlerFunc()
	//
	// mux.HandleFunc("/api/lobby/players", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("got here")
	//
	// 	ws, err := upgrader.Upgrade(w, r, nil)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	defer ws.Close()
	//
	// 	if err := ws.WriteMessage(websocket.TextMessage, []byte("sup")); err != nil {
	// 		fmt.Println("sup", err)
	// 	}
	// 	if err := ws.WriteMessage(websocket.TextMessage, []byte("bob")); err != nil {
	// 		fmt.Println("sup", err)
	// 	}
	// 	if err := ws.WriteMessage(websocket.TextMessage, []byte("jon")); err != nil {
	// 		fmt.Println("sup", err)
	// 	}
	// })
}
