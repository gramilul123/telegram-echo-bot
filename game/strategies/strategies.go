package strategies

import (
	"fmt"
	"math/rand"
	"time"

	war_map "github.com/gramilul123/telegram-echo-bot/game/war_map"
)

const (
	DONE      = "done"
	HIT       = "hit"
	NOK       = "nok"
	WIN       = "win"
	DESTROYED = "destroyed"
)

const (
	SHIP  = 3
	HALO  = 4
	EMPTY = 5
)

const (
	SIMPLE = "simple"
	MIDDLE = "middle"
)

type Strategy interface {
	GetShot(result string) (int, int, [][]int)
	Create()
}

func GetStrategy(variant string) (strategy Strategy) {

	switch variant {
	case MIDDLE:
		strategy = &MiddleStrategy{}
	default:
		strategy = &SimpleStrategy{}
	}

	return
}

func CheckShot(x, y int, gameWorkMap [][]int, gameMap war_map.WarMap) (string, [][]int) {
	var result string

	if x == 0 && y == 0 {

		return DONE, gameWorkMap
	}

	if gameMap.Cells[x][y] == war_map.Ship {
		gameWorkMap[x][y] = SHIP

		allShipsDestroed := true
		for i, ship := range gameMap.WarShips {

			if ship.Status {

				index := fmt.Sprintf("%v-%v", x, y)
				if _, ok := ship.Location[index]; ok {
					delete(gameMap.WarShips[i].Location, index)

					if len(gameMap.WarShips[i].Location) == 0 {
						for _, coordinate := range gameMap.WarShips[i].HaloLocation {
							gameWorkMap[coordinate.X][coordinate.Y] = HALO
						}
						gameMap.WarShips[i].Status = false
						result = DESTROYED
					} else {
						result = HIT
					}
				}
			}
		}

		for _, ship := range gameMap.WarShips {
			if ship.Status {
				allShipsDestroed = false
			}
		}

		if allShipsDestroed {
			return WIN, gameWorkMap
		}

	} else {
		gameWorkMap[x][y] = EMPTY
		result = NOK
	}

	return result, gameWorkMap
}

func GetRndVariant(variants []war_map.Coordinate) (x, y int) {

	switch len(variants) {
	case 0:
		x, y = 0, 0
	case 1:
		x, y = variants[0].X, variants[0].Y
	default:
		source := rand.NewSource(time.Now().UnixNano())
		rand := rand.New(source)

		rndKey := rand.Intn(len(variants))
		x, y = variants[rndKey].X, variants[rndKey].Y
	}

	return
}

func SimplyGetingShot(gameMap [][]int) (x, y int) {
	var variants []war_map.Coordinate

	for i := 1; i <= 10; i++ {
		for y := 1; y <= 10; y++ {
			if gameMap[i][y] == 0 {
				variants = append(variants, war_map.Coordinate{i, y})
			}
		}
	}

	x, y = GetRndVariant(variants)

	return x, y
}
