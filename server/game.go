package server

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var games = make(map[string]*Game)
var gameMu sync.Mutex

type GamePlayer struct {
	Player Player
	ws     *websocket.Conn
}

type Game struct {
	Id      string
	gameMu  sync.Mutex
	players []GamePlayer
}

func NewGame() *Game {
	uuid := uuid.New()

	g := Game{
		Id: uuid.String(),
	}
	gameMu.Lock()
	games[g.Id] = &g
	gameMu.Unlock()

	return &g
}

func StartGameFromLobby(players []Player) *Game {
	gp := make([]GamePlayer, len(players))
	for i, p := range players {
		gp[i] = GamePlayer{
			Player: p,
			ws:     nil,
		}
	}

	uuid := uuid.New()
	g := Game{
		Id:      uuid.String(),
		players: gp,
	}

	gameMu.Lock()
	games[g.Id] = &g
	gameMu.Unlock()

	return &g
}

func GetGame(gid string) (*Game, bool) {
	gameMu.Lock()
	defer gameMu.Unlock()
	g, exists := games[gid]
	return g, exists
}

func (g *Game) GetPlayer(playerId string) (*GamePlayer, bool) {
	g.gameMu.Lock()
	defer g.gameMu.Unlock()

	for _, p := range g.players {
		if p.Player.Id == playerId {
			return &p, true
		}
	}
	return nil, false

}

func (g *Game) GetPlayers() []GamePlayer {
	g.gameMu.Lock()
	defer g.gameMu.Unlock()

	return g.players
}
