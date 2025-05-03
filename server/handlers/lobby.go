package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/OndrejHana/two-towers/server/lobby"
	"github.com/clerk/clerk-sdk-go/v2"
)

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

type GetLobbyResponse struct {
	Players []lobby.Player `json:"players"`
	Ready   bool           `json:"ready"`
}

type StatusResponse struct {
	Event   lobby.Event    `json:"event"`
	Players []lobby.Player `json:"players"`
}

type ToggleReadyResponse struct {
	Player lobby.Player `json:"player"`
}

func NewLobbyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionClaims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: No session found", http.StatusUnauthorized)
		return
	}

	clerkUserID := sessionClaims.Subject

	var req CreateLobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	l, err := lobby.NewLobby(4)
	if err != nil {
		http.Error(w, "Lobby could not be created", http.StatusInternalServerError)
		return
	}

	p := lobby.NewPlayer(clerkUserID, req.Username)
	if err := l.AddPlayer(p); err != nil {
		http.Error(w, "Player could not be added", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateLobbyResponse{Code: l.GetCode()})
}

func GetLobbyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionClaims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: No session found", http.StatusUnauthorized)
		return
	}

	clerkUserID := sessionClaims.Subject

	code := r.URL.Query().Get(":code")
	l, exists := lobby.GetLobby(code)

	if !exists {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	p, exists := l.GetPlayer(clerkUserID)
	if !exists {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetLobbyResponse{
		Players: l.GetPlayers(),
		Ready:   p.Ready,
	})
}

func LobbyStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionClaims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: No session found", http.StatusUnauthorized)
		return
	}

	clerkUserID := sessionClaims.Subject

	code := r.URL.Query().Get(":code")
	l, exists := lobby.GetLobby(code)

	if !exists {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	p, exists := l.GetPlayer(clerkUserID)
	fmt.Println(clerkUserID, p, exists)
	if !exists {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	var res StatusResponse
	select {
	case e := <-p.GetEventChan():
		res = StatusResponse{
			Event:   e,
			Players: l.GetPlayers(),
		}
	case <-time.After(time.Second * 5):
		res = StatusResponse{
			Event: lobby.Event{
				EventType: lobby.Status,
			},
			Players: l.GetPlayers(),
		}
	}

	json.NewEncoder(w).Encode(res)

}

func JoinLobbyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JoinLobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	sessionClaims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: No session found", http.StatusUnauthorized)
		return
	}

	clerkUserID := sessionClaims.Subject

	l, exists := lobby.GetLobby(req.Code)

	if !exists {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	p := lobby.NewPlayer(clerkUserID, req.Username)
	if err := l.AddPlayer(p); err != nil {
		http.Error(w, "Player could not be added", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LeaveLobbyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionClaims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: No session found", http.StatusUnauthorized)
		return
	}

	clerkUserID := sessionClaims.Subject

	code := r.URL.Query().Get(":code")
	l, exists := lobby.GetLobby(code)
	if !exists {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	_, deleted := l.RemovePlayer(clerkUserID)
	if deleted {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "User not in lobby", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
	}
}

func ToggleReadyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionClaims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: No session found", http.StatusUnauthorized)
		return
	}

	clerkUserID := sessionClaims.Subject

	code := r.URL.Query().Get(":code")
	l, exists := lobby.GetLobby(code)
	if !exists {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	_, err := l.ToggleReady(clerkUserID)
	if err != nil {
		http.Error(w, "Could not toggle ready state", http.StatusNotFound)
		return
	}

	p, exists := l.GetPlayer(clerkUserID)
	if !exists {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToggleReadyResponse{Player: *p})
}
