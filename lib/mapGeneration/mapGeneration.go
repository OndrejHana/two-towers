package mapgeneration

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"sort"
)

type Structure uint8

const (
	NONE Structure = iota
	TOWER
	ROAD
)

type Field struct {
	structure Structure
	id        uint
	p         Point
}

type Point struct {
	x, y uint
}

func getDistance(a, b Point) float64 {
	x1 := float64(a.x)
	x2 := float64(b.x)

	y1 := float64(a.y)
	y2 := float64(b.y)

	dx := math.Abs(x1 - x2)
	dy := math.Abs(y1 - y2)

	return dx + dy
}

func createGrid(partsPerSide uint, fieldsPerTower uint) [][]Field {
	side := (partsPerSide * fieldsPerTower) + (partsPerSide + 1)
	var grid [][]Field = make([][]Field, side, side)
	for row := range grid {
		grid[row] = make([]Field, side, side)
	}

	return grid
}

func printGrid(grid [][]Field) {
	for i := range grid {
		fmt.Printf("%d: %v\n", i, grid[i])
	}
}

func createTowers(grid [][]Field, v uint, partsPerSide uint, fieldsPerTower uint) map[uint]Point {
	towers := make(map[uint]Point)
	var id uint = 0
	for i := range partsPerSide {
		for j := range partsPerSide {
			starty := i*fieldsPerTower + i + 1
			startx := j*fieldsPerTower + j + 1

			num := uint(rand.Intn(int((fieldsPerTower * fieldsPerTower) + 1)))
			deltay := num / fieldsPerTower
			deltax := num % fieldsPerTower

			y := starty + deltay
			x := startx + deltax
			p := Point{
				x: x,
				y: y,
			}
			t := Field{
				structure: TOWER,
				id:        id,
			}

			grid[y][x] = t
			towers[id] = p

			id += 1
			if id >= v {
				break
			}
		}
	}

	return towers
}

type Edge struct {
	first    uint
	second   uint
	distance uint
}

func calculateDistanceBetweenTowers(towers map[uint]Point) []Edge {
	list := []Edge{}

	for t1, p1 := range towers {
		for t2 := t1 + 1; t2 < uint(len(towers)); t2++ {
			p2 := towers[t2]

			distance := uint(math.Round(getDistance(p1, p2)))
			list = append(list, Edge{t1, t2, distance})
		}
	}

	return list
}

type DisjointSet struct {
	Parent map[Field]Field
	Rank   map[Field]int
}

func NewDisjointSet(points []Field) *DisjointSet {
	ds := &DisjointSet{
		Parent: make(map[Field]Field),
		Rank:   make(map[Field]int),
	}
	for _, p := range points {
		ds.Parent[p] = p
		ds.Rank[p] = 0
	}
	return ds
}

func (ds *DisjointSet) Find(p Field) Field {
	if ds.Parent[p] == p {
		return p
	}
	root := ds.Find(ds.Parent[p])
	ds.Parent[p] = root
	return root
}

func (ds *DisjointSet) Union(a, b Field) {
	rootA := ds.Find(a)
	rootB := ds.Find(b)

	if rootA != rootB {
		if ds.Rank[rootA] < ds.Rank[rootB] {
			ds.Parent[rootA] = rootB
		} else if ds.Rank[rootA] > ds.Rank[rootB] {
			ds.Parent[rootB] = rootA
		} else {
			ds.Parent[rootB] = rootA
			ds.Rank[rootA]++
		}
	}
}

func getMST(paths []Edge, towers []Field) []Edge {
	ds := NewDisjointSet(towers)

	mst := []Edge{}

	for _, edge := range paths {
		if ds.Find(towers[edge.first]) != ds.Find(towers[edge.second]) {
			mst = append(mst, edge)
			ds.Union(towers[edge.first], towers[edge.second])
		}
	}

	return mst
}

func isInGrid(p Point, grid [][]Field) bool {
	side := uint(len(grid))
	return p.x >= 0 && p.x < side && p.y >= 0 && p.y < side
}

func isRoad(p Point, grid [][]Field) bool {
	return grid[p.x][p.y].structure == ROAD
}

type AStarNode struct {
	Point    Point
	Priority float64
	GScore   int
}

// PriorityQueue is a priority queue for A* search.
type PriorityQueue []*AStarNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(*AStarNode)
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}

func findPath(start, end Point, grid [][]Field) []Point {
	heuristic := func(a, b Point) float64 {
		return getDistance(a, b)
	}

	getNeighbors := func(p Point) []Field {
		neighbors := []Point{
			{x: p.x, y: p.y + 1},
			{x: p.x + 1, y: p.y},
			{x: p.x, y: p.y - 1},
			{x: p.x - 1, y: p.y},
		}

		validNeighbors := []Field{}
		for _, neighbor := range neighbors {
			if isInGrid(neighbor, grid) && !isRoad(neighbor, grid) {
				gridVal := grid[neighbor.x][neighbor.y]
				gridVal.p = neighbor
				validNeighbors = append(validNeighbors, gridVal)
			}
		}

		return validNeighbors
	}

	openSet := &PriorityQueue{&AStarNode{Point: start, Priority: 0, GScore: 0}}
	heap.Init(openSet)

	cameFrom := make(map[Point]Point)
	gScore := make(map[Point]int)
	gScore[start] = 0

	fScore := make(map[Point]float64)
	fScore[start] = heuristic(start, end)

	for openSet.Len() > 0 {
		currentNode := heap.Pop(openSet).(*AStarNode)
		currentPoint := currentNode.Point

		if currentPoint == end {
			// Reconstruct path
			path := []Point{currentPoint}
			for {
				prev, ok := cameFrom[currentPoint]
				if !ok {
					break
				}
				path = append(path, prev)
				currentPoint = prev
			}
			// Reverse the path
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}
			return path
		}

		for _, neighbor := range getNeighbors(currentPoint) {
			tempGScore := gScore[currentPoint] + 1

			if _, ok := gScore[neighbor.p]; !ok || tempGScore < gScore[neighbor.p] {
				cameFrom[neighbor.p] = currentPoint
				gScore[neighbor.p] = tempGScore
				fScore[neighbor.p] = float64(tempGScore) + (rand.Float64()*2 - 1) + heuristic(neighbor.p, end)

				heap.Push(openSet, &AStarNode{Point: neighbor.p, Priority: fScore[neighbor.p], GScore: tempGScore})
			}
		}
	}

	return nil // No path found
}

func GenerateMap(n uint, fieldsPerTower uint) {
	v := 3*n + 2
	partsPerSide := uint(math.Ceil(math.Sqrt(float64(v))))
	grid := createGrid(partsPerSide, fieldsPerTower)
	printGrid(grid)
	fmt.Println()

	towers := createTowers(grid, v, partsPerSide, fieldsPerTower)

	printGrid(grid)
	fmt.Println()
	fmt.Printf("towers: %v\n", towers)
	fmt.Println()

	distances := calculateDistanceBetweenTowers(towers)
	fmt.Printf("distances: %v\n", distances)
	fmt.Println()

	sort.Slice(distances, func(i, j int) bool {
		d1 := distances[i]
		d2 := distances[j]

		return d1.distance < d2.distance
	})

	fmt.Printf("distances: %v\n", distances)
	fmt.Println()

	towerList := []Field{}
	for key := range towers {
		towerList = append(towerList, Field{
			structure: TOWER,
			id:        key,
		})
	}

	paths := getMST(distances, towerList)

	fmt.Printf("paths: %v\n", paths)

	for _, p := range paths {
		fmt.Println(findPath(towers[p.first], towers[p.second], grid))
	}
}
