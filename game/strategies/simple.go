package strategies

import (
	"encoding/json"
	"log"

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

func (s SimpleStrategy) MapToJson() string {
	str, err := json.Marshal(s)

	if err != nil {
		log.Fatal(err)
	}

	return (string(str))
}

func (s *SimpleStrategy) JsonToMap(str string) {
	bytes := []byte(str)

	if err := json.Unmarshal(bytes, &s); err != nil {
		log.Fatal(err)
	}
}
