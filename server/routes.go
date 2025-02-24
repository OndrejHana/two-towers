package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"

	"two-towers/lib/lobbyStore"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

func getAuth(r *http.Request) (*clerk.User, bool, error) {
	ctx := r.Context()
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		return nil, false, nil
	}

	u, err := user.Get(ctx, claims.Subject)
	if err != nil {
		return nil, ok, err
	}

	return u, true, nil
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

	mux.Handle("/api/lobby/new", clerkhttp.WithHeaderAuthorization()(newLobbyHandler))
}
