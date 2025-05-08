package game

import (
	"fmt"
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

type MessagePlayer struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Color    string `json:"color"`
}

type Point struct {
	X uint `json:"x"`
	Y uint `json:"y"`
}

type Tile struct {
	Structure   Structure `json:"structure"`
	StructureId *uint     `json:"structureId"`
	UnitId      *uint     `json:"unitId"`
}

type Unit struct {
	PlayerId uint  `json:"playerId"`
	Point    Point `json:"point"`
}

type Tower struct {
	PlayerId     *uint `json:"playerId"`
	TargetRoadId *uint `json:"targetRoadId"`
	Point        Point `json:"point"`
}

type Road struct {
	Points []Point `json:"points"`
	From   Tower   `json:"from"`
	To     Tower   `json:"to"`
}

type World struct {
	Grid  [][]Tile
	Roads []Road
}

type State struct {
	Towers []Tower
	Units  []Unit
}

func InitializeGrid(width, height int) [][]Tile {
	grid := make([][]Tile, height)
	for i := range grid {
		grid[i] = make([]Tile, width)
		for j := range grid[i] {
			grid[i][j] = Tile{
				Structure:   NONE,
				StructureId: nil,
				UnitId:      nil,
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

// func placeTowers(grid [][]Tile, n int) []Tower {
// 	towers := []Tower{}
// 	height := len(grid)
// 	width := len(grid[0])
// 	for len(towers) < n {
// 		y := uint(rand.Intn(height))
// 		x := uint(rand.Intn(width))
// 		for grid[y][x].Structure != NONE {
// 			y = uint(rand.Intn(height))
// 			x = uint(rand.Intn(width))
// 		}
// 		towers = append(towers, Tower{
// 			x:        x,
// 			y:        y,
// 			PlayerId: nil,
// 		})
// 		tid := uint(len(towers) - 1)
// 		grid[y][x] = Tile{
// 			Structure:   TOWER,
// 			StructureId: &tid,
// 			UnitId:      nil,
// 		}
// 	}
// 	return towers
// }
//
// func isValid(p Point, grid [][]Tile) bool {
// 	return p.x >= 0 && p.x < uint(len(grid)) && p.y >= 0 && p.y < uint(len(grid[0]))
// }
//
// func dfs(grid [][]Tile, current, target Point, visited map[Point]bool, path []Point) []Point {
// 	if current == target {
// 		return path
// 	}
//
// 	for _, d := range directions {
// 		next := Point{current.x + uint(d[0]), current.y + uint(d[1])}
// 		if !isValid(next, grid) {
// 			continue
// 		}
// 		// Avoid crossing an already drawn path (unless this is the target cell)
// 		if grid[next.y][next.x].Structure == ROAD && next != target {
// 			continue
// 		}
// 		// Avoid revisiting the same cell
// 		if visited[next] {
// 			continue
// 		}
// 		// Allow stepping only into an empty cell, or into a tower if it is the target
// 		if grid[next.y][next.x].Structure != NONE &&
// 			!(grid[next.y][next.x].Structure == TOWER && next == target) {
// 			continue
// 		}
// 		visited[next] = true
// 		newPath := append(path, next)
// 		if result := dfs(grid, next, target, visited, newPath); result != nil {
// 			return result
// 		}
// 		delete(visited, next)
// 	}
// 	return nil
// }
//
// func ConnectTowers(grid [][]Tile, tower1, tower2 Tower) bool {
// 	visited := make(map[Point]bool)
// 	visited[tower1.getPoint()] = true
// 	path := dfs(grid, tower1.getPoint(), tower2.getPoint(), visited, []Point{tower1.getPoint()})
// 	if path == nil {
// 		fmt.Printf("No path found between %v and %v\n", tower1, tower2)
// 		return false
// 	}
// 	// Mark the path cells; ensure not to override tower cells.
// 	for _, cell := range path {
// 		if grid[cell.y][cell.x].Structure != TOWER {
// 			grid[cell.y][cell.x].Structure = ROAD
// 		}
// 	}
// 	return true
// }

func CreateMock(players []MessagePlayer) (World, State) {
	units := []Unit{}

	towers := []Tower{
		{
			PlayerId: nil,
			Point: Point{
				X: 1,
				Y: 1,
			},
		},
		{
			PlayerId: nil,
			Point: Point{
				X: 8,
				Y: 1,
			},
		},
		{
			PlayerId: nil,
			Point: Point{
				X: 1,
				Y: 8,
			},
		},
		{
			PlayerId: nil,
			Point: Point{
				X: 8,
				Y: 8,
			},
		},
	}

	width, height := 10, 10
	world := InitializeGrid(width, height)

	// Place towers
	t := uint(0)
	redTower := &world[1][1]
	redTower.Structure = TOWER
	redTower.StructureId = &t

	t2 := uint(1)
	blueTower := &world[1][8]
	blueTower.Structure = TOWER
	blueTower.StructureId = &t2

	t3 := uint(2)
	greenTower := &world[8][1]
	greenTower.Structure = TOWER
	greenTower.StructureId = &t3

	t4 := uint(3)
	purpleTower := &world[8][8]
	purpleTower.Structure = TOWER
	purpleTower.StructureId = &t4

	roads := []Road{}

	// Horizontal road 1 (top)
	tiles1 := []Point{}
	for x := 2; x < 8; x++ {
		if world[1][x].Structure == 0 {
			world[1][x].Structure = ROAD
			world[1][x].StructureId = &t
			tiles1 = append(tiles1, Point{X: uint(x), Y: 1})
		}
	}
	roads = append(roads, Road{
		Points: tiles1,
		From:   towers[0],
		To:     towers[1],
	})

	// Horizontal road 2 (bottom)
	tiles2 := []Point{}
	for x := 2; x < 8; x++ {
		if world[8][x].Structure == 0 {
			world[8][x].Structure = ROAD
			world[8][x].StructureId = &t2
			tiles2 = append(tiles2, Point{X: uint(x), Y: 8})
		}
	}
	roads = append(roads, Road{
		Points: tiles2,
		From:   towers[2],
		To:     towers[3],
	})

	// Vertical road 1 (left)
	tiles3 := []Point{}
	for y := 2; y < 8; y++ {
		if world[y][1].Structure == 0 {
			world[y][1].Structure = ROAD
			world[y][1].StructureId = &t3
			tiles3 = append(tiles3, Point{X: 1, Y: uint(y)})
		}
	}
	roads = append(roads, Road{
		Points: tiles3,
		From:   towers[0],
		To:     towers[2],
	})

	// Vertical road 2 (right)
	tiles4 := []Point{}
	for y := 2; y < 8; y++ {
		if world[y][8].Structure == 0 {
			world[y][8].Structure = ROAD
			world[y][8].StructureId = &t4
			tiles4 = append(tiles4, Point{X: 8, Y: uint(y)})
		}
	}
	roads = append(roads, Road{
		Points: tiles4,
		From:   towers[1],
		To:     towers[3],
	})

	return World{
			Grid:  world,
			Roads: roads,
		}, State{
			Towers: towers,
			Units:  units,
		}
}
