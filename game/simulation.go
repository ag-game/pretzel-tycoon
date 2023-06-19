package game

import (
	"math/rand"
)

type SimulationActor struct {
	X, Y             int
	TargetX, TargetY int
	TargetActive     bool
	WalkTicks        int
	WaitTicks        int
	Inactive         bool
}

type Simulation struct {
	Money  int
	Actors []*SimulationActor
}

func (s *Simulation) generateActors() {
	for i := 0; i < 7; i++ {
		x := -1
		tx := 41
		if rand.Intn(2) == 0 {
			x = 41
			tx = -1
		}
		y := -1
		ty := 18
		if rand.Intn(2) == 0 {
			y = 18
			ty = -1
		}
		a := &SimulationActor{
			X:         x,
			Y:         y,
			TargetX:   tx,
			TargetY:   ty,
			WalkTicks: 20 + rand.Intn(30),
			WaitTicks: rand.Intn(500),
		}
		s.Actors = append(s.Actors, a)
	}
}

func (s *Simulation) StartDay() {
	s.Actors = s.Actors[:]
	s.generateActors()
}

func (s *Simulation) Tick() error {
	standX, standY := 19, 9

	var tx, ty int
	for _, a := range s.Actors {
		if a.Inactive {
			continue
		}

		if a.WaitTicks > 0 {
			a.WaitTicks--
			continue
		}

		if a.TargetActive {
			tx, ty = a.TargetX, a.TargetY
		} else {
			tx, ty = standX, standY
		}

		if a.X < tx {
			a.X++
		} else if a.X > tx {
			a.X--
		}
		if a.Y < ty {
			a.Y++
		} else if a.Y > ty {
			a.Y--
		}

		if a.X == tx && a.Y == ty {
			if !a.TargetActive {
				a.TargetActive = true
				a.WaitTicks = 100 + rand.Intn(200)
				continue
			} else {
				a.Inactive = true
				continue
			}
		}
		a.WaitTicks = a.WalkTicks
	}
	return nil
}
