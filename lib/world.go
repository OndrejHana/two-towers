package lib

import (
	"fmt"
	"math/rand"
)

var directions = [][2]int{
	{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
}

type Structure uint

const (
	NONE Structure = iota
	TOWER
	ROAD
)

type Point struct {
	x uint
	y uint
}

type Tile struct {
	Structure Structure `json:"structure"`
	TowerId   *uint     `json:"towerId"`
	UnitId    *uint     `json:"unitId"`
}

type Unit struct {
	PlayerId uint `json:"playerId"`
}

type Tower struct {
	PlayerId *uint `json:"playerId"`
	x        uint
	y        uint
}

func (t *Tower) getPoint() Point {
	return Point{
		x: t.x,
		y: t.y,
	}
}

type Player struct {
	Color string `json:"color"`
}

type Game struct {
	World   [][]Tile `json:"world"`
	Players []Player `json:"players"`
	Towers  []Tower  `json:"towers"`
	Units   []Unit   `json:"units"`
}

func InitializeGrid(width, height int) [][]Tile {
	grid := make([][]Tile, height)
	for i := range grid {
		grid[i] = make([]Tile, width)
		for j := range grid[i] {
			grid[i][j] = Tile{
				Structure: NONE,
				TowerId:   nil,
				UnitId:    nil,
			}
		}
	}
	return grid
}

func PrintGrid(grid [][]Tile) {
	for _, row := range grid {
		for _, cell := range row {
			fmt.Print(cell.Structure, " ")
		}
		fmt.Println()
	}
	fmt.Println()
}

func placeTowers(grid [][]Tile, n int) []Tower {
	towers := []Tower{}
	height := len(grid)
	width := len(grid[0])
	for len(towers) < n {
		y := uint(rand.Intn(height))
		x := uint(rand.Intn(width))
		for grid[y][x].Structure != NONE {
			y = uint(rand.Intn(height))
			x = uint(rand.Intn(width))
		}
		towers = append(towers, Tower{
			x:        x,
			y:        y,
			PlayerId: nil,
		})
		tid := uint(len(towers) - 1)
		grid[y][x] = Tile{
			Structure: TOWER,
			TowerId:   &tid,
			UnitId:    nil,
		}
	}
	return towers
}

func isValid(p Point, grid [][]Tile) bool {
	return p.x >= 0 && p.x < uint(len(grid)) && p.y >= 0 && p.y < uint(len(grid[0]))
}

func dfs(grid [][]Tile, current, target Point, visited map[Point]bool, path []Point) []Point {
	if current == target {
		return path
	}

	for _, d := range directions {
		next := Point{current.x + uint(d[0]), current.y + uint(d[1])}
		if !isValid(next, grid) {
			continue
		}
		// Avoid crossing an already drawn path (unless this is the target cell)
		if grid[next.y][next.x].Structure == ROAD && next != target {
			continue
		}
		// Avoid revisiting the same cell
		if visited[next] {
			continue
		}
		// Allow stepping only into an empty cell, or into a tower if it is the target
		if grid[next.y][next.x].Structure != NONE &&
			!(grid[next.y][next.x].Structure == TOWER && next == target) {
			continue
		}
		visited[next] = true
		newPath := append(path, next)
		if result := dfs(grid, next, target, visited, newPath); result != nil {
			return result
		}
		delete(visited, next)
	}
	return nil
}

func ConnectTowers(grid [][]Tile, tower1, tower2 Tower) bool {
	visited := make(map[Point]bool)
	visited[tower1.getPoint()] = true
	path := dfs(grid, tower1.getPoint(), tower2.getPoint(), visited, []Point{tower1.getPoint()})
	if path == nil {
		fmt.Printf("No path found between %v and %v\n", tower1, tower2)
		return false
	}
	// Mark the path cells; ensure not to override tower cells.
	for _, cell := range path {
		if grid[cell.y][cell.x].Structure != TOWER {
			grid[cell.y][cell.x].Structure = ROAD
		}
	}
	return true
}

func CreateMock() Game {
	players := []Player{
		{Color: "#0000ff"},
		{Color: "#00ff00"},
		{Color: "#ff0000"},
	}

	units := []Unit{
		{PlayerId: 1},
	}

	p1 := uint(0)
	p2 := uint(1)
	towers := []Tower{
		{PlayerId: &p1},
		{PlayerId: nil},
		{PlayerId: &p2},
	}

	width, height := 10, 10
	world := InitializeGrid(width, height)

	t := uint(0)
	redTower := &world[1][1]
	redTower.Structure = TOWER
	redTower.TowerId = &t

	t2 := uint(1)
	blueTower := &world[8][8]
	blueTower.Structure = TOWER
	blueTower.TowerId = &t2

	t3 := uint(2)
	greenTower := &world[1][8]
	greenTower.Structure = TOWER
	greenTower.TowerId = &t3

	for x := 2; x < 8; x++ {
		if world[1][x].Structure == 0 {
			world[1][x].Structure = ROAD
		}
	}
	for y := 2; y < 8; y++ {
		if world[y][8].Structure == 0 {
			world[y][8].Structure = ROAD
		}
	}

	world[1][4].UnitId = &t

	return Game{
		World:   world,
		Players: players,
		Units:   units,
		Towers:  towers,
	}
}
