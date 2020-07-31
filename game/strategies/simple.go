package strategies

import (
	war_map "github.com/gramilul123/telegram-echo-bot/game/war_map"
)

type SimpleStrategy struct {
	GameWorkMap [][]int
}

func (s *SimpleStrategy) Create() {
	gameWorkMap := war_map.WarMap{}
	gameWorkMap.Create(false)

	s.GameWorkMap = gameWorkMap.Cells
}

func (s SimpleStrategy) GetShot(result string) (int, int, [][]int) {

	x, y := SimplyGetingShot(s.GameWorkMap)

	return x, y, s.GameWorkMap
}
