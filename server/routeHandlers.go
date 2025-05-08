package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
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
	Players []Player `json:"players"`
	Ready   bool     `json:"ready"`
}

type StatusResponse struct {
	Event   Event    `json:"event"`
	Players []Player `json:"players"`
}

type ToggleReadyResponse struct {
	Player Player `json:"player"`
}

type WSUserId struct {
	UserToken string `json:"userToken"`
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

	l, err := NewLobby(4)
	if err != nil {
		http.Error(w, "Lobby could not be created", http.StatusInternalServerError)
		return
	}

	p := NewPlayer(clerkUserID, req.Username)
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
	l, exists := GetLobby(code)

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
	l, exists := GetLobby(code)

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
			Event: Event{
				EventType: Status,
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

	l, exists := GetLobby(req.Code)

	if !exists {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	p := NewPlayer(clerkUserID, req.Username)
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
	l, exists := GetLobby(code)
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
	l, exists := GetLobby(code)
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

func HandleGameWs(w http.ResponseWriter, r *http.Request) {
	gid := r.URL.Query().Get(":gameId")
	g, exists := GetGame(gid)
	if !exists {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	t, payload, err := ws.ReadMessage()
	if err != nil {
		http.Error(w, "Could not read payload", http.StatusNotFound)
		return
	}

	fmt.Println(t)

	var uid WSUserId
	if err := json.Unmarshal(payload, &uid); err != nil {
		http.Error(w, "Invalid content", http.StatusNotFound)
		return
	}

	claims, err := jwt.Verify(r.Context(), &jwt.VerifyParams{
		Token: uid.UserToken,
	})

	if err != nil {
		http.Error(w, "Could not get user", http.StatusInternalServerError)
		return
	}

	usr, err := user.Get(r.Context(), claims.Subject)
	if err != nil {
		http.Error(w, "Could not get user", http.StatusInternalServerError)
		return
	}

	if exists = g.ConnectPlayer(usr.ID, ws); !exists {
		http.Error(w, "Player not found in game", http.StatusInternalServerError)
		return
	}
}
