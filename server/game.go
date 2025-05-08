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
	MessageType int        `json:"messageType"`
	World       game.World `json:"world"`
	State       game.State `json:"state"`
	Players     []Player   `json:"players"`
}

type GamePlayer struct {
	Player Player
	Color  string
	ws     *websocket.Conn
}

type Game struct {
	Id             string
	gameMu         sync.Mutex
	players        []GamePlayer
	playerConnChan chan string
	world          game.World
	state          game.State
}

func StartGameFromLobby(players []Player) *Game {
	gp := make([]GamePlayer, len(players))
	colors := GeneratePlayerColorStrings(len(players))
	for i, p := range players {
		gp[i] = GamePlayer{
			Player: p,
			Color:  colors[i],
			ws:     nil,
		}
	}

	uuid := uuid.New()

	world := game.GenerateWorld(len(players))
	fmt.Println(world)

	// Get player IDs for state generation
	playerIds := make([]string, len(players))
	for i, p := range players {
		playerIds[i] = p.Id
	}

	// Generate initial state
	state := game.GenerateInitialState(world, playerIds)

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
				MessageType: 1, // Initial game state
				World:       g.world,
				State:       g.state,
				Players:     players,
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

func (g *Game) StartGameLoop() error {
	return nil
}
