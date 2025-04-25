package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OndrejHana/two-towers/server/lobby"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
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
	router.Handle("/lobby/new", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreateLobbyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		newLobby, err := lobby.NewLobby()
		if err != nil {
			http.Error(w, "Failed to create lobby", http.StatusInternalServerError)
			return
		}

		// Add the creator as the first player
		player := lobby.Player{
			ID:       r.Header.Get("Authorization"),
			Username: req.Username,
		}
		if err := newLobby.AddPlayer(player); err != nil {
			http.Error(w, "Failed to add player to lobby", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateLobbyResponse{Code: newLobby.Code})
	})))

	router.Handle("/lobby/join", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JoinLobbyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		l, exists := lobby.GetLobby(req.Code)
		if !exists {
			http.Error(w, "Lobby not found", http.StatusNotFound)
			return
		}

		player := lobby.Player{
			ID:       r.Header.Get("Authorization"),
			Username: req.Username,
		}

		if err := l.AddPlayer(player); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"players": l.GetPlayers(),
		})
	}))

	router.Handle("/lobby/{code}/players", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get(":code")
		if code == "" {
			// Try to get the code from the URL path
			code = r.URL.Path[len("/lobby/"):]
			code = code[:len(code)-len("/players")]
		}

		l, exists := lobby.GetLobby(code)
		if !exists {
			http.Error(w, "Lobby not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"players": l.GetPlayers(),
		})
	})))

	router.Handle("/lobby/{code}/ready", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ready called")
		code := r.URL.Query().Get(":code")
		if code == "" {
			// Try to get the code from the URL path
			code = r.URL.Path[len("/lobby/"):]
			code = code[:len(code)-len("/ready")]
		}

		l, exists := lobby.GetLobby(code)
		if !exists {
			http.Error(w, "Lobby not found", http.StatusNotFound)
			return
		}

		playerID := r.Header.Get("Authorization")
		if err := l.ToggleReady(playerID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if all players are ready
		allReady := true
		players := l.GetPlayers()
		for _, p := range players {
			fmt.Println(p.Username, p.Ready)
			if !p.Ready {
				allReady = false
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"gameStarted": allReady && len(players) > 1,
		})
	})))

	router.Handle("/game/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("Authorization")
		if token == "" {
			http.Error(w, "Missing session token", http.StatusUnauthorized)
			return
		}

		if _, err := jwt.Verify(r.Context(), &jwt.VerifyParams{
			Token: token,
			JWK:   nil,
		}); err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("upgrade:", err)
			return
		}
		defer c.Close()

		// TODO: Handle game websocket messages
	}))

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("dist"))).Methods("GET")
}
