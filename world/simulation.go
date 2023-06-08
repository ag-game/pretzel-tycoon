package world

import "github.com/hajimehoshi/ebiten/v2"

// SimulationItems are all in grams.
type SimulationItems struct {
	Flour  int
	Water  int
	Salt   int
	Yeast  int
	Sugar  int
	Butter int
}

type SimulationActor struct {
	X, Y  float64
	Image *ebiten.Image
}

type Simulation struct {
	Money int

	Mix      *SimulationItems
	Supplies *SimulationItems

	Actors []*SimulationActor
}

func (s *Simulation) Tick() error {
	return nil
}
