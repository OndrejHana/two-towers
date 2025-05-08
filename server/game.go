package server

import (
	"fmt"
	"sync"

	"github.com/OndrejHana/two-towers/server/game"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var games = make(map[string]*Game)
var gameMu sync.Mutex

type Message struct {
	World   game.World           `json:"world"`
	Players []game.MessagePlayer `json:"players"`
}

type GamePlayer struct {
	Player        Player
	MessagePlayer game.MessagePlayer
	ws            *websocket.Conn
	color         string
}

type Game struct {
	Id             string
	gameMu         sync.Mutex
	players        []GamePlayer
	playerConnChan chan string
	world          game.World
	state          game.State
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
	colors := GeneratePlayerColorStrings(len(players))
	for i, p := range players {
		gp[i] = GamePlayer{
			Player: p,
			MessagePlayer: game.MessagePlayer{
				Id:       p.Id,
				Username: p.Username,
				Color:    colors[i],
			},
			ws:    nil,
			color: colors[i],
		}
	}

	uuid := uuid.New()

	mp := make([]game.MessagePlayer, len(players))
	for i, gp := range gp {
		mp[i] = gp.MessagePlayer
	}

	world, state := game.CreateMock(mp)

	g := Game{
		Id:             uuid.String(),
		players:        gp,
		playerConnChan: make(chan string, len(players)),
		world:          world,
		state:          state,
	}

	gameMu.Lock()
	games[g.Id] = &g
	gameMu.Unlock()

	go func() {
		for range g.playerConnChan {
			allConnected := true
			for _, p := range g.players {
				if p.ws == nil {
					allConnected = false
				}
			}
			fmt.Println(g.players, allConnected)
			if !allConnected {
				continue
			}

			players := make([]Player, len(g.players))
			for i, p := range g.players {
				players[i] = p.Player
			}

			// All players are connected, send initial game state
			message := Message{
				World:   g.world,
				Players: mp,
			}

			for _, p := range g.players {
				if err := p.ws.WriteJSON(message); err != nil {
					fmt.Printf("Error sending game state to player %s: %v\n", p.Player.Id, err)
				}
			}

			break
		}

		fmt.Println("Game started, initial state sent to all players")
	}()

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

func (g *Game) ConnectPlayer(playerId string, ws *websocket.Conn) bool {
	gameMu.Lock()
	defer gameMu.Unlock()

	for i, p := range g.players {
		if p.Player.Id == playerId {
			if p.ws != nil {
				return false
			}
			g.players[i].ws = ws
			g.playerConnChan <- playerId
			return true
		}
	}

	return false

}
