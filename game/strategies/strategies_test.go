package strategies

import (
	"testing"

	war_map "github.com/gramilul123/telegram-echo-bot/game/war_map"
)

func TestSimpleStrategy(t *testing.T) {
	var result string

	gameMap := war_map.WarMap{}
	gameMap.Create(true)

	str := GetStrategy(SIMPLE)
	str.Create()

	x, y, gameWorkMap := str.GetShot(result)
	result, _ = CheckShot(x, y, gameWorkMap, gameMap)

	hasError := false
	switch gameMap.Cells[x][y] {
	case war_map.Ship:
		if !(result == HIT || result == DESTROYED) {
			hasError = true
		}
	default:
		if result != NOK {
			hasError = true
		}
	}

	if hasError {
		t.Errorf("Error simple strategy")
	}
}

func TestMiddleStrategy(t *testing.T) {
	var result string
	var x, y int
	var gameWorkMap [][]int

	gameMap := war_map.WarMap{}
	gameMap.Create(true)

	variantStrategy := GetStrategy(MIDDLE)
	variantStrategy.Create()

	itt := 1
	for {
		x, y, gameWorkMap = variantStrategy.GetShot(result)
		result, _ = CheckShot(x, y, gameWorkMap, gameMap)

		if result == WIN || result == DONE {
			break
		}
		itt++
	}

	if itt > 100 || !allShipsDestroyed(gameWorkMap) {
		t.Errorf("Error middle strategy")
	}
}

func allShipsDestroyed(gameWorkMap [][]int) bool {
	count := 0

	for i, row := range gameWorkMap {
		if i > 0 && i < 11 {
			for j, cell := range row {
				if j > 0 && j < 11 {
					if cell == 3 {
						count++
					}
				}
			}
		}
	}

	return count == 20
}
