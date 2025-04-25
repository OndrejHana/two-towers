package lobby

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
)

type Player struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Ready    bool   `json:"ready"`
}

type Lobby struct {
	Code    string
	Players []Player
	mu      sync.Mutex
	MaxSize int
}

var lobbies = make(map[string]*Lobby)
var mu sync.Mutex

func NewLobby() (*Lobby, error) {
	code, err := generateCode()
	if err != nil {
		return nil, err
	}

	lobby := &Lobby{
		Code:    code,
		Players: make([]Player, 0),
		MaxSize: 4,
	}

	mu.Lock()
	lobbies[code] = lobby
	mu.Unlock()

	return lobby, nil
}

func GetLobby(code string) (*Lobby, bool) {
	mu.Lock()
	defer mu.Unlock()
	lobby, exists := lobbies[code]
	return lobby, exists
}

func (l *Lobby) AddPlayer(player Player) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.Players) >= l.MaxSize {
		return fmt.Errorf("lobby is full")
	}

	l.Players = append(l.Players, player)
	fmt.Println(l.Players)
	return nil
}

func (l *Lobby) RemovePlayer(playerID string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i, p := range l.Players {
		if p.ID == playerID {
			l.Players = append(l.Players[:i], l.Players[i+1:]...)
			fmt.Println(l.Players)
			break
		}
	}
}

func (l *Lobby) GetPlayers() []Player {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.Players
}

func (l *Lobby) ToggleReady(playerID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i, p := range l.Players {
		if p.ID == playerID {
			l.Players[i].Ready = !l.Players[i].Ready
			fmt.Println(p.Username, l.Players[i].Ready)
			return nil
		}
	}
	return fmt.Errorf("player not found")
}

func generateCode() (string, error) {
	bytes := make([]byte, 2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:4], nil
}
