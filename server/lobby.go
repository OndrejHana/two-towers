package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"slices"
	"sync"
)

type EventType int

const (
	AddPlayer = iota
	RemovePlayer
	ToggleReady
	Status
	Start
)

func generateCode() (string, error) {
	bytes := make([]byte, 2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:4], nil
}

type Event struct {
	EventType int
	GameId    *string
}

type Player struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Ready     bool   `json:"ready"`
	eventChan chan Event
}

func NewPlayer(id string, username string) Player {
	return Player{
		Id:        id,
		Username:  username,
		eventChan: make(chan Event, 16),
	}
}

func (p *Player) GetEventChan() <-chan Event {
	return p.eventChan
}

type Lobby struct {
	code    string
	players []Player
	events  chan Event
	lobbyMu sync.Mutex
	maxSize int
}

var lobbies = make(map[string]*Lobby)
var lobbyMu sync.Mutex

func (l *Lobby) GetCode() string {
	return l.code
}

func NewLobby(maxSize int) (*Lobby, error) {
	code, err := generateCode()
	if err != nil {
		return nil, err
	}

	lobby := Lobby{
		code:    code,
		players: []Player{},
		events:  make(chan Event, maxSize*4),
		maxSize: maxSize,
	}

	go func() {
		for e := range lobby.events {
			allReady := true

			for _, p := range lobby.players {
				pe := Event{EventType: e.EventType}
				if e.EventType == Start {
					g := StartGameFromLobby(lobby.players)
					pe.GameId = &g.Id
				}
				p.eventChan <- pe

				if !p.Ready {
					allReady = false
				}
			}

			if e.EventType == ToggleReady && allReady {
				lobby.events <- Event{EventType: Start}
			}
		}
	}()

	lobbyMu.Lock()
	lobbies[code] = &lobby
	lobbyMu.Unlock()

	return &lobby, nil
}

func (l *Lobby) AddPlayer(player Player) error {
	l.lobbyMu.Lock()
	defer l.lobbyMu.Unlock()

	if len(l.players) >= l.maxSize {
		return errors.New("lobby is full")
	}
	l.players = append(l.players, player)
	l.events <- Event{
		EventType: AddPlayer,
	}

	return nil
}

func GetLobby(code string) (*Lobby, bool) {
	lobbyMu.Lock()
	defer lobbyMu.Unlock()
	lobby, exists := lobbies[code]
	return lobby, exists
}

func (l *Lobby) GetPlayers() []Player {
	l.lobbyMu.Lock()
	defer l.lobbyMu.Unlock()
	return l.players
}

func (l *Lobby) GetPlayer(playerId string) (*Player, bool) {
	l.lobbyMu.Lock()
	defer l.lobbyMu.Unlock()

	for _, p := range l.players {
		if p.Id == playerId {
			return &p, true
		}
	}
	return nil, false
}

func (l *Lobby) ToggleReady(playerID string) (bool, error) {
	l.lobbyMu.Lock()
	defer l.lobbyMu.Unlock()

	for i, p := range l.players {
		if p.Id == playerID {
			l.players[i].Ready = !l.players[i].Ready

			l.events <- Event{
				EventType: ToggleReady,
			}

			return l.players[i].Ready, nil
		}
	}
	return false, errors.New("player not found")
}

func (l *Lobby) RemovePlayer(playerID string) (*Player, bool) {
	l.lobbyMu.Lock()
	defer l.lobbyMu.Unlock()

	for i, p := range l.players {
		if p.Id == playerID {
			player := l.players[i]
			l.players = slices.Delete(l.players, i, i+1)

			l.events <- Event{
				EventType: RemovePlayer,
			}

			return &player, true
		}
	}

	return nil, false
}
