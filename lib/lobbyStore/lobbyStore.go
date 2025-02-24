package lobbystore

import (
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
	id             uuid.UUID
	connString     string
	players        uint8
	p1, p2, p3, p4 *string
}

func (l *Lobby) GetConnString() string {
	return l.connString
}

func NewLobby(p1 string) Lobby {
	id := uuid.Must(uuid.NewV7())
	connString := randSeq(4)

	return Lobby{
		id:         id,
		connString: connString,
		players:    1,
		p1:         &p1,
		p2:         nil,
		p3:         nil,
		p4:         nil,
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
