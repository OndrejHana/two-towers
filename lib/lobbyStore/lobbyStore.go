package lobbystore

import (
	"encoding/json"
	"math/rand"

	"github.com/google/uuid"
)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type NewLobbyRes struct {
	ConnString string
}

type Lobby struct {
	id          uuid.UUID
	connString  string
	playerCount uint8
	players     []string
}

func (l *Lobby) GetConnString() string {
	return l.connString
}

func (l *Lobby) JSON() ([]byte, error) {
	return json.Marshal(l)
}

func NewLobby(p1 string) Lobby {
	id := uuid.Must(uuid.NewV7())
	connString := randSeq(4)

	return Lobby{
		id:          id,
		connString:  connString,
		playerCount: 1,
		players:     []string{p1},
	}
}

type LobbyStore struct {
	lobbies map[string]Lobby
}

func NewLobbyStore() LobbyStore {
	return LobbyStore{
		lobbies: make(map[string]Lobby),
	}
}

func (ls *LobbyStore) NewLobby(userId string) *Lobby {
	l := NewLobby(userId)
	ls.lobbies[l.connString] = l
	return &l
}
