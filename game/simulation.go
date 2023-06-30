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
	Day int

	Money            int
	Actors           []*SimulationActor
	DayFinished      bool
	DayFinishedTicks int
	PretzelsSold     int

	MakePretzels int
	MakeSigns    int

	PretzelPrice int // In cents.
}

func NewSimulation() *Simulation {
	return &Simulation{
		Day: 1,

		MakePretzels: -1,
		MakeSigns:    -1,
		PretzelPrice: -1,
	}
}

func (s *Simulation) generateActors() {
	for i := 0; i < 7; i++ {
		x := -2 + rand.Intn(40)
		tx := 44 - x
		y := 18
		ty := 18 + rand.Intn(15)
		a := &SimulationActor{
			X:         x,
			Y:         y,
			TargetX:   tx,
			TargetY:   ty,
			WalkTicks: 10 + rand.Intn(30),
			WaitTicks: rand.Intn(500),
		}
		s.Actors = append(s.Actors, a)
	}
}

func (s *Simulation) StartDay() {
	s.Actors = s.Actors[:0]
	s.DayFinished = false
	s.DayFinishedTicks = 25
	s.PretzelsSold = 0

	s.generateActors()
}

func (s *Simulation) Tick() error {
	if s.DayFinished {
		return nil
	}

	standX, standY := 19, 10

	var active int
	var tx, ty int
	for _, a := range s.Actors {
		if a.Inactive {
			continue
		}
		active++

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

		if a.TargetActive && (a.X < -2 || a.X > 40 || a.Y > 17) {
			a.Inactive = true
			continue
		}

		if a.X == tx && a.Y == ty {
			if !a.TargetActive {
				s.PretzelsSold++

				a.TargetActive = true
				a.WaitTicks = 100 + rand.Intn(150)
				continue
			} else {
				a.Inactive = true
				continue
			}
		}
		a.WaitTicks = a.WalkTicks
	}
	if active == 0 {
		s.DayFinishedTicks--
		if s.DayFinishedTicks == 0 {
			s.DayFinished = true
		}
	}
	return nil
}
