package lib

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	grid := InitializeGrid(10, 10)
	towers := placeTowers(grid, 4)
	fmt.Println(towers)
	PrintGrid(grid)
	fmt.Println(ConnectTowers(grid, towers[0], towers[1]))
	PrintGrid(grid)
}
