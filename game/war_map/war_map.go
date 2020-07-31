package war_map

import (
	"fmt"
	"math/rand"
	"time"

	warship "github.com/gramilul123/telegram-echo-bot/game/warship"
)

const (
	SizeMap = 12
)

const (
	Empty = 0
	Ship  = 1
	Halo  = 2

	Left   = 1
	Top    = 2
	Right  = 3
	Bottom = 4
)

var CountShipsRule = map[int]int{
	warship.Battleship:  1,
	warship.Cruiser:     2,
	warship.Destroyer:   3,
	warship.TorpedoBoat: 4,
}

type DirectionVariant struct {
	X, Y, Direction, ShipType int
}

var Directions = []DirectionVariant{}

type WarMap struct {
	WarShips []warship.Warship
	Cells    [][]int
}

type Coordinate struct {
	X, Y int
}

var HaloCoordinates = []warship.HaloLocation{}

func (mapData *WarMap) Create(addShip bool) {

	cellsMap := make([][]int, SizeMap)
	for i := range cellsMap {
		cellsMap[i] = make([]int, SizeMap)
	}
	mapData.Cells = cellsMap

	if addShip {
		for shipType, count := range CountShipsRule {
			for i := 0; i < count; i++ {
				created, location, haloLocation := AddShipIfPossible(shipType, mapData.Cells)

				if created {
					mapData.WarShips = append(mapData.WarShips, warship.Warship{
						Type:         shipType,
						Location:     location,
						HaloLocation: haloLocation,
						Status:       true,
					})
				}
			}
		}
	}
}

func AddShipIfPossible(shipType int, cells [][]int) (bool, map[string]bool, []warship.HaloLocation) {
	var location map[string]bool
	var haloLocation []warship.HaloLocation

	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)

	for {
		x, y := rand.Intn(SizeMap-2)+1, rand.Intn(SizeMap-2)+1

		Directions = CheckCoordinates(cells, shipType, x, y)

		if len(Directions) > 0 {
			break
		}
	}

	if len(Directions) > 0 {

		var variant int
		if len(Directions) > 1 {
			variant = rand.Intn(len(Directions))
		}

		location, haloLocation = AddCoordinates(cells, Directions[variant])

		return true, location, haloLocation
	}

	return false, location, haloLocation
}

func CheckCoordinates(cells [][]int, shipType, x, y int) []DirectionVariant {
	Directions = []DirectionVariant{}

	if cells[x][y] == Empty {
		if y-shipType+1 >= 1 {
			CheckCoordinate(&Directions, cells, shipType, x, y, Left)
		}

		if x-shipType+1 >= 1 {
			CheckCoordinate(&Directions, cells, shipType, x, y, Top)
		}

		if y+shipType-1 <= SizeMap-2 {
			CheckCoordinate(&Directions, cells, shipType, x, y, Right)
		}

		if x+shipType-1 <= SizeMap-2 {
			CheckCoordinate(&Directions, cells, shipType, x, y, Bottom)
		}
	}

	return Directions
}

func CheckCoordinate(Directions *[]DirectionVariant, cells [][]int, shipType, x, y, direction int) {
	valid := true

	for i := 0; i < shipType; i++ {

		if direction == Left {
			if cells[x][y-i] != Empty {
				valid = false
			}
		} else if direction == Top {
			if cells[x-i][y] != Empty {
				valid = false
			}
		} else if direction == Right {
			if cells[x][y+i] != Empty {
				valid = false
			}
		} else if direction == Bottom {
			if cells[x+i][y] != Empty {
				valid = false
			}
		}

		if !valid {
			break
		}
	}

	if valid {
		*Directions = append(*Directions, DirectionVariant{x, y, direction, shipType})
	}
}

func AddCoordinates(cells [][]int, direction DirectionVariant) (map[string]bool, []warship.HaloLocation) {

	var dx, dy int
	location := make(map[string]bool)
	HaloCoordinates = []warship.HaloLocation{}

	for i := 0; i < direction.ShipType; i++ {

		if direction.Direction == Left {
			dx = direction.X
			dy = direction.Y - i

			HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + 1, Y: direction.Y - i})
			HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - 1, Y: direction.Y - i})

			if i == 0 {
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X, Y: direction.Y - i + 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + 1, Y: direction.Y - i + 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - 1, Y: direction.Y - i + 1})
			}

			if i == direction.ShipType-1 {
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X, Y: direction.Y - i - 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + 1, Y: direction.Y - i - 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - 1, Y: direction.Y - i - 1})
			}

		} else if direction.Direction == Top {
			dx = direction.X - i
			dy = direction.Y

			HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - i, Y: direction.Y + 1})
			HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - i, Y: direction.Y - 1})

			if i == 0 {
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - i + 1, Y: direction.Y})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - i + 1, Y: direction.Y + 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - i + 1, Y: direction.Y - 1})
			}

			if i == direction.ShipType-1 {
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - i - 1, Y: direction.Y})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - i - 1, Y: direction.Y + 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - i - 1, Y: direction.Y - 1})
			}

		} else if direction.Direction == Right {
			dx = direction.X
			dy = direction.Y + i

			HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + 1, Y: direction.Y + i})
			HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - 1, Y: direction.Y + i})

			if i == 0 {
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X, Y: direction.Y + i - 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + 1, Y: direction.Y + i - 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - 1, Y: direction.Y + i - 1})
			}

			if i == direction.ShipType-1 {
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X, Y: direction.Y + i + 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + 1, Y: direction.Y + i + 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X - 1, Y: direction.Y + i + 1})
			}

		} else if direction.Direction == Bottom {
			dx = direction.X + i
			dy = direction.Y

			HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + i, Y: direction.Y + 1})
			HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + i, Y: direction.Y - 1})

			if i == 0 {
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + i - 1, Y: direction.Y})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + i - 1, Y: direction.Y + 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + i - 1, Y: direction.Y - 1})
			}

			if i == direction.ShipType-1 {
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + i + 1, Y: direction.Y})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + i + 1, Y: direction.Y + 1})
				HaloCoordinates = append(HaloCoordinates, warship.HaloLocation{X: direction.X + i + 1, Y: direction.Y - 1})
			}
		}

		cells[dx][dy] = Ship
		for _, coordinates := range HaloCoordinates {
			cells[coordinates.X][coordinates.Y] = Halo
		}

		index := fmt.Sprintf("%v-%v", dx, dy)
		location[index] = true
	}

	return location, HaloCoordinates
}
