package giorengine

import (
	"encoding/json"
	"io"
	"math"
)

type Match struct {
	Name          string
	Width, Height int
	Users         []FactionSpawn
	Notes         string
	TurnCount     int
	Cities        []CitySpawn
	Mountains     []int
	Moves         []Move
	AFKs          [][]int
}

func (m *Match) UnmarshalJSON(data []byte) error {
	users := []string{}
	stars := []float64{}
	moves := [][]int{}
	cities := []int{}
	cityarmies := []int{}
	generals := []int{}
	array := []interface{}{
		nil,
		&m.Name,
		&m.Width,
		&m.Height,
		&users,
		&stars,
		&cities,
		&cityarmies,
		&generals,
		&m.Mountains,
		&moves,
		&m.AFKs,
	}

	err := json.Unmarshal(data, &array)

	for i, user := range users {
		m.Users = append(m.Users, struct {
			Pos   int
			Name  string
			Stars float64
		}{generals[i], user, stars[i]})
	}

	for i, city := range cities {
		m.Cities = append(m.Cities, struct {
			Pos   int
			Count int
		}{city, cityarmies[i]})
	}

	for _, move := range moves {
		m.Moves = append(m.Moves, Move{
			move[0], move[1], move[2], move[3] == 1, move[4],
		})
	}

	m.TurnCount = int(math.Ceil(float64(moves[len(moves)-1][4]) / 2))

	return err
}

func ParseReplay(r io.Reader) (*Match, error) {
	der, _ := NewReaderUint16BE(r)
	decoder := json.NewDecoder(der)

	var m *Match
	err := decoder.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func LoadReplay(r io.Reader) (*Game, error) {
	m, err := ParseReplay(r)
	if err != nil {
		return nil, err
	}
	return New(m), nil
}
