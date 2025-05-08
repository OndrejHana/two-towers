package game

import (
	"math"
	"math/rand"
)

var directions = [][2]int{
	{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
}

func getTowerCount(n int) int {
	return 3*n + 2
}

// GenerateWorld creates a new game world with randomly distributed towers and roads
func GenerateWorld(numPlayers int) World {
	// Calculate grid size based on number of towers
	numTowers := getTowerCount(numPlayers)
	gridSize := numTowers * 2 // Make grid large enough to space towers

	// Initialize empty grid
	grid := make([][]Tile, gridSize)
	for i := range grid {
		grid[i] = make([]Tile, gridSize)
	}

	// Place towers randomly
	towers := make([]Tower, 0, numTowers)
	towerId := uint(1)
	for len(towers) < numTowers {
		x := uint(rand.Intn(gridSize))
		y := uint(rand.Intn(gridSize))

		// Check if position is valid (not too close to other towers)
		valid := true
		for _, tower := range towers {
			dx := int(x) - int(tower.Point.X)
			dy := int(y) - int(tower.Point.Y)
			if dx*dx+dy*dy < 4 { // Minimum distance of 2 tiles
				valid = false
				break
			}
		}

		if valid {
			// Place tower
			grid[y][x].StructureType = StructureTypeTower
			grid[y][x].StructureId = &towerId
			towers = append(towers, Tower{
				Id:       towerId,
				PlayerId: "", // Will be assigned later
				Point:    Point{X: x, Y: y},
			})
			towerId++
		}
	}

	// Connect towers with roads
	roads := make([]Road, 0)
	roadId := uint(1)
	connections := make(map[uint]int) // Track number of connections per tower

	// First, ensure each tower has at least one connection
	for i := 0; i < len(towers); i++ {
		if connections[towers[i].Id] == 0 {
			// Find closest unconnected tower
			bestDist := -1
			bestTower := -1
			for j := 0; j < len(towers); j++ {
				if i != j && connections[towers[j].Id] < 3 {
					dx := int(towers[i].Point.X) - int(towers[j].Point.X)
					dy := int(towers[i].Point.Y) - int(towers[j].Point.Y)
					dist := dx*dx + dy*dy
					if bestDist == -1 || dist < bestDist {
						bestDist = dist
						bestTower = j
					}
				}
			}
			if bestTower != -1 {
				// Create road between towers
				road := createRoad(roadId, towers[i], towers[bestTower], grid)
				roads = append(roads, road)
				connections[towers[i].Id]++
				connections[towers[bestTower].Id]++
				roadId++
			}
		}
	}

	// Add additional random connections
	for i := 0; i < len(towers); i++ {
		for connections[towers[i].Id] < 3 {
			// Find random tower to connect to
			bestDist := -1
			bestTower := -1
			for j := 0; j < len(towers); j++ {
				if i != j && connections[towers[j].Id] < 3 {
					dx := int(towers[i].Point.X) - int(towers[j].Point.X)
					dy := int(towers[i].Point.Y) - int(towers[j].Point.Y)
					dist := dx*dx + dy*dy
					if bestDist == -1 || dist < bestDist {
						bestDist = dist
						bestTower = j
					}
				}
			}
			if bestTower != -1 {
				// Create road between towers
				road := createRoad(roadId, towers[i], towers[bestTower], grid)
				roads = append(roads, road)
				connections[towers[i].Id]++
				connections[towers[bestTower].Id]++
				roadId++
			} else {
				break
			}
		}
	}

	return World{
		Grid:  grid,
		Roads: roads,
	}
}

// createRoad creates a road between two towers and updates the grid
func createRoad(id uint, from, to Tower, grid [][]Tile) Road {
	points := make([]Point, 0)

	// Simple straight line path for now
	x1, y1 := int(from.Point.X), int(from.Point.Y)
	x2, y2 := int(to.Point.X), int(to.Point.Y)

	// Add points along the path
	points = append(points, from.Point)

	// Add intermediate points
	if x1 != x2 {
		dx := 1
		if x2 < x1 {
			dx = -1
		}
		for x := x1 + dx; x != x2; x += dx {
			points = append(points, Point{X: uint(x), Y: uint(y1)})
			grid[y1][x].StructureType = StructureTypeRoad
			grid[y1][x].StructureId = &id
		}
	}

	if y1 != y2 {
		dy := 1
		if y2 < y1 {
			dy = -1
		}
		for y := y1 + dy; y != y2; y += dy {
			points = append(points, Point{X: uint(x2), Y: uint(y)})
			grid[y][x2].StructureType = StructureTypeRoad
			grid[y][x2].StructureId = &id
		}
	}

	points = append(points, to.Point)

	return Road{
		Id:     id,
		Points: points,
		From:   from,
		To:     to,
	}
}

// GenerateInitialState creates the initial game state with one tower assigned to each player
func GenerateInitialState(world World, playerIds []string) State {
	state := State{
		Towers: make([]Tower, 0),
		Units:  make([]Unit, 0),
	}

	// Calculate grid dimensions
	gridWidth := len(world.Grid[0])
	gridHeight := len(world.Grid)

	// Calculate sector size based on number of players
	sectorsX := int(math.Ceil(math.Sqrt(float64(len(playerIds)))))
	sectorsY := int(math.Ceil(float64(len(playerIds)) / float64(sectorsX)))

	sectorWidth := gridWidth / sectorsX
	sectorHeight := gridHeight / sectorsY

	// Create a map to track which towers have been assigned
	assignedTowers := make(map[uint]bool)

	// Assign towers to players, one per sector
	for i, playerId := range playerIds {
		// Calculate sector coordinates
		sectorX := i % sectorsX
		sectorY := i / sectorsX

		// Calculate sector boundaries
		startX := sectorX * sectorWidth
		endX := (sectorX + 1) * sectorWidth
		startY := sectorY * sectorHeight
		endY := (sectorY + 1) * sectorHeight

		// Find the best tower in this sector
		bestTower := findBestTowerInSector(world, startX, endX, startY, endY, assignedTowers)
		if bestTower != nil {
			bestTower.PlayerId = playerId
			state.Towers = append(state.Towers, *bestTower)
			assignedTowers[bestTower.Id] = true
		}
	}

	return state
}

// findBestTowerInSector finds the best tower in a given sector
func findBestTowerInSector(world World, startX, endX, startY, endY int, assignedTowers map[uint]bool) *Tower {
	var bestTower *Tower
	centerX := (startX + endX) / 2
	centerY := (startY + endY) / 2
	minDist := -1

	// Look for towers in the sector
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			if world.Grid[y][x].StructureType == StructureTypeTower && world.Grid[y][x].StructureId != nil {
				towerId := *world.Grid[y][x].StructureId
				if !assignedTowers[towerId] {
					// Calculate distance to sector center
					dx := x - centerX
					dy := y - centerY
					dist := dx*dx + dy*dy

					// Prefer towers closer to the center of the sector
					if minDist == -1 || dist < minDist {
						minDist = dist
						bestTower = &Tower{
							Id:       towerId,
							PlayerId: "", // Will be set by caller
							Point:    Point{X: uint(x), Y: uint(y)},
						}
					}
				}
			}
		}
	}

	return bestTower
}
