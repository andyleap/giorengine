// giorengine project giorengine.go
package giorengine

import (
	"math"
)

type Game struct {
	Cells    []*Cell
	Factions []*Faction
	Turn     int
	Match    *Match
	moves    []Move
	afks     [][]int
	afkPos   int
}

type Cell struct {
	Type    CellType
	Pop     int
	Faction *Faction
}

type CellType int

const (
	Plain CellType = iota
	City
	Mountain
	General
)

type Faction struct {
	Name   string
	Pop    int
	Land   int
	Cities int
	AFK    bool
	ID     int
}

func (g *Game) PreTurn() {
	g.Turn++
	GenCities := false
	GenAll := false
	if g.Turn%2 == 0 {
		GenCities = true
	}
	if g.Turn%50 == 0 {
		GenAll = true
	}
	for _, c := range g.Cells {
		if GenCities && (c.Type == City || c.Type == General) {
			if c.Faction != nil || c.Pop < 50 {
				c.Pop++
				if c.Faction != nil {
					c.Faction.Pop++
				}
			}
		}
		if GenAll && c.Faction != nil {
			c.Pop++
			c.Faction.Pop++
		}
	}
}

type CitySpawn struct {
	Pos   int
	Count int
}

type FactionSpawn struct {
	Pos   int
	Name  string
	Stars float64
}

func New(m *Match) *Game {
	g := &Game{}
	g.Match = m
	g.Cells = make([]*Cell, m.Width*m.Height)
	for i := range g.Cells {
		g.Cells[i] = &Cell{}
	}
	for _, city := range m.Cities {
		g.Cells[city.Pos].Type = City
		g.Cells[city.Pos].Pop = city.Count
	}
	for _, mountain := range m.Mountains {
		g.Cells[mountain].Type = Mountain
	}
	for i, faction := range m.Users {
		f := &Faction{Name: faction.Name, Pop: 1, Land: 1, Cities: 1, ID: i}
		g.Cells[faction.Pos].Type = General
		g.Cells[faction.Pos].Faction = f
		g.Cells[faction.Pos].Pop = 1
		g.Factions = append(g.Factions, f)
	}
	g.moves = m.Moves
	g.afks = m.AFKs
	g.afkPos = 0
	return g
}

func (g *Game) Move(player int, from, to int, is50 bool) {
	fc, tc := g.Cells[from], g.Cells[to]
	if fc.Faction == nil {
		return
	}

	amount := 1
	if is50 {
		amount = int(math.Ceil(float64(fc.Pop) / 2))
	}
	amount = fc.Pop - amount
	fc.Pop -= amount

	if tc.Faction != fc.Faction {
		fc.Faction.Pop -= amount
		if tc.Faction != nil {
			tc.Faction.Pop -= amount
		}
		tc.Pop -= amount
		if tc.Pop < 0 {
			fc.Faction.Pop -= tc.Pop
			fc.Faction.Land++
			if tc.Faction != nil {
				tc.Faction.Pop -= tc.Pop
				tc.Faction.Land--
			}
			if tc.Type == City || tc.Type == General {
				if tc.Faction != nil {
					tc.Faction.Cities--
				}
				fc.Faction.Cities++
			}
			oldFac := tc.Faction
			tc.Faction = fc.Faction
			tc.Pop = -tc.Pop

			if tc.Type == General {
				tc.Type = City
				for _, cell := range g.Cells {
					if cell.Faction == oldFac {
						cell.Faction = fc.Faction
						oldFac.Pop -= cell.Pop
						oldFac.Land--
						cell.Pop = int(math.Floor(float64(cell.Pop)/2 + 0.5))
						fc.Faction.Land++
						fc.Faction.Pop += cell.Pop
					}
				}
			}
		}
	} else {
		tc.Pop += amount
	}
}

type Move struct {
	Player   int
	From, To int
	Is50     bool
	Turn     int
}

func (g *Game) Step() bool {
	g.PreTurn()
	for g.afkPos < len(g.afks) && g.Turn == g.afks[g.afkPos][1] {
		f := g.Factions[g.afks[g.afkPos][0]]
		if f.AFK {
			for _, c := range g.Cells {
				if c.Faction == f {
					c.Faction = nil
					if c.Type == General {
						c.Type = City
					}
				}
			}
			f.Land = 0
			f.Pop = 0
		}
		f.AFK = true
		g.afkPos++
	}
	action := len(g.moves) > 0
	for len(g.moves) > 0 && g.moves[0].Turn <= g.Turn {
		move := g.moves[0]
		g.Move(move.Player, move.From, move.To, move.Is50)
		g.moves = g.moves[1:]
	}
	return action
}
