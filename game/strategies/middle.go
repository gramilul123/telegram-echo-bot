package strategies

import (
	"encoding/json"
	"log"

	war_map "github.com/gramilul123/telegram-echo-bot/game/war_map"
)

type MiddleStrategy struct {
	GameWorkMap  [][]int
	LastShot     war_map.Coordinate
	VariantsHits map[string]war_map.Coordinate
}

func (s *MiddleStrategy) Create() {
	gameWorkMap := war_map.WarMap{}
	gameWorkMap.Create(false)

	s.GameWorkMap = gameWorkMap.Cells
	s.VariantsHits = make(map[string]war_map.Coordinate)
}

func (s *MiddleStrategy) GetShot(result string) (int, int, [][]int) {
	var x, y int
	var variantsHits []war_map.Coordinate
	var onlyHorizontal, onlyVertical bool

	if result == NOK && len(s.VariantsHits) > 0 {

		for key, value := range s.VariantsHits {
			if s.LastShot.X == value.X && s.LastShot.Y == value.Y {
				delete(s.VariantsHits, key)
			}
		}

		for _, value := range s.VariantsHits {
			variantsHits = append(variantsHits, value)
		}

		x, y = GetRndVariant(variantsHits)

	} else if result == HIT {

		for key, value := range s.VariantsHits {
			if s.LastShot.X == value.X && s.LastShot.Y == value.Y {
				delete(s.VariantsHits, key)
			}
		}

		if s.LastShot.Y-1 > 0 {
			if s.GameWorkMap[s.LastShot.X][s.LastShot.Y-1] == 0 {
				s.VariantsHits["left"] = war_map.Coordinate{s.LastShot.X, s.LastShot.Y - 1}
			} else if s.GameWorkMap[s.LastShot.X][s.LastShot.Y-1] == 3 {
				onlyHorizontal = true
			}
		}
		if s.LastShot.Y+1 <= 10 {
			if s.GameWorkMap[s.LastShot.X][s.LastShot.Y+1] == 0 {
				s.VariantsHits["right"] = war_map.Coordinate{s.LastShot.X, s.LastShot.Y + 1}
			} else if s.GameWorkMap[s.LastShot.X][s.LastShot.Y+1] == 3 {
				onlyHorizontal = true
			}
		}
		if s.LastShot.X-1 > 0 {
			if s.GameWorkMap[s.LastShot.X-1][s.LastShot.Y] == 0 {
				s.VariantsHits["top"] = war_map.Coordinate{s.LastShot.X - 1, s.LastShot.Y}
			} else if s.GameWorkMap[s.LastShot.X-1][s.LastShot.Y] == 3 {
				onlyVertical = true
			}
		}
		if s.LastShot.X+1 <= 10 {
			if s.GameWorkMap[s.LastShot.X+1][s.LastShot.Y] == 0 {
				s.VariantsHits["bottom"] = war_map.Coordinate{s.LastShot.X + 1, s.LastShot.Y}
			} else if s.GameWorkMap[s.LastShot.X+1][s.LastShot.Y] == 3 {
				onlyVertical = true
			}
		}

		if onlyHorizontal {
			delete(s.VariantsHits, "top")
			delete(s.VariantsHits, "bottom")
		} else if onlyVertical {
			delete(s.VariantsHits, "left")
			delete(s.VariantsHits, "right")
		}

		for _, value := range s.VariantsHits {
			variantsHits = append(variantsHits, value)
		}

		x, y = GetRndVariant(variantsHits)
	}

	if (result == NOK && len(s.VariantsHits) == 0) || result == DESTROYED || result == "" || (x == 0 && y == 0) {

		for key, _ := range s.VariantsHits {
			delete(s.VariantsHits, key)
		}

		x, y = SimplyGetingShot(s.GameWorkMap)
	}
	s.LastShot = war_map.Coordinate{x, y}

	return x, y, s.GameWorkMap
}

func (s MiddleStrategy) MapToJson() string {
	str, err := json.Marshal(s)

	if err != nil {
		log.Fatal(err)
	}

	return (string(str))
}

func (s *MiddleStrategy) JsonToMap(str string) {
	bytes := []byte(str)

	if err := json.Unmarshal(bytes, &s); err != nil {
		log.Fatal(err)
	}
}
