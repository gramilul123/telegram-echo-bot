package warship

const (
	Battleship  = 4
	Cruiser     = 3
	Destroyer   = 2
	TorpedoBoat = 1
)

type HaloLocation struct {
	X int
	Y int
}

type Warship struct {
	Type         int
	Location     map[string]bool
	HaloLocation []HaloLocation
	Status       bool
}
