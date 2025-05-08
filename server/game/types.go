package game

type StructureType uint

const (
	StructureTypeNone StructureType = iota
	StructureTypeTower
	StructureTypeRoad
)

type Point struct {
	X uint `json:"x"`
	Y uint `json:"y"`
}

type Tile struct {
	StructureType StructureType `json:"structureType"`
	StructureId   *uint         `json:"structureId"`
}

type Tower struct {
	Id       uint   `json:"id"`
	PlayerId string `json:"playerId"`
	Point    Point  `json:"point"`
}

type Unit struct {
	Id            uint   `json:"id"`
	PlayerId      string `json:"playerId"`
	Point         Point  `json:"point"`
	RoadId        uint   `json:"roadId"`
	TargetTowerId uint   `json:"targetTowerId"`
}

type Road struct {
	Id     uint    `json:"id"`
	Points []Point `json:"points"`
	From   Tower   `json:"from"`
	To     Tower   `json:"to"`
}

type World struct {
	Grid  [][]Tile `json:"grid"`
	Roads []Road   `json:"roads"`
}

type State struct {
	Towers []Tower `json:"towers"`
	Units  []Unit  `json:"units"`
}
