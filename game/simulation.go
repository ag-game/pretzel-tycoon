package game

import (
	"math/rand"
)

const (
	PretzelCost = 15
	SignCost    = 125
)

const (
	DisasterWeather = iota + 1
	DisasterParentsLostJobs
	DisasterMomMedical
	DisasterDadMedical
	DisasterPlayerMedical
	DisasterCompetitor
	DisasterChildLabor
	DisasterBankAccount
	DisasterCreditCard
	DisasterInflation
	DisasterDeflation
	DisasterStockMarket
	DisasterThugs
	DisasterDebt
	DisasterSued
	DisasterWar
)

func disasterLabel(d int) string {
	switch d {
	case DisasterWeather:
		return "RAIN, RAIN, GO AWAY!"
	case DisasterParentsLostJobs:
		return "YOUR PARENTS LOST BOTH OF THEIR JOBS!"
	case DisasterMomMedical:
		return "YOUR MOM NEEDS EXPENSIVE MEDICAL WORK!"
	case DisasterDadMedical:
		return "YOUR DAD NEEDS EXPENSIVE MEDICAL WORK!"
	case DisasterPlayerMedical:
		return "YOU NEED EXPENSIVE MEDICAL WORK!"
	case DisasterCompetitor:
		return "A MEGA-CORP COMPETITOR OPENS NEARBY!"
	case DisasterChildLabor:
		return "YOUR ARE ACCUSED OF USING CHILD LABOR!"
	case DisasterBankAccount:
		return "YOUR BANK ACCOUNT IS FROZEN!"
	case DisasterCreditCard:
		return "YOUR CREDIT CARD IS FROZEN!"
	case DisasterInflation:
		return "MASSIVE INFLATION PLAGUES THE ECONOMY!"
	case DisasterDeflation:
		return "MASSIVE DEFLATION PLAGUES THE ECONOMY!"
	case DisasterStockMarket:
		return "THE STOCK MARKET HAS CRASHED!"
	case DisasterThugs:
		return "THUGS KNOCK OVER YOUR PRETZEL STAND!"
	case DisasterDebt:
		return "YOUR CAR IS REPOSSESSED DUE TO DEBT!"
	case DisasterSued:
		return "A CUSTOMER IS SUING YOU FOR FRAUD!"
	case DisasterWar:
		return "WAR HAS BEEN DECLARED!"
	default:
		return ""
	}
}

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

	NoSaleNoStock int

	Disaster    int
	DisasterDay int // Last disaster day.
	Disasters   []int

	DisasterAcknowledged bool
}

func NewSimulation() *Simulation {
	return &Simulation{
		Day:   1,
		Money: 420,

		MakePretzels: -1,
		MakeSigns:    -1,
		PretzelPrice: -1,
	}
}

func (s *Simulation) nextDisaster() int {
	if s.Day == 1 || len(s.Disasters) == 16 {
		return 0
	}

	if s.Day-s.DisasterDay < 3 && rand.Intn(2) == 0 {
		return 0
	}

DISASTERS:
	for {
		disaster := rand.Intn(16) + 1
		for _, d := range s.Disasters {
			if d == disaster {
				continue DISASTERS
			}
		}
		s.Disasters = append(s.Disasters, disaster)
		s.DisasterDay = s.Day
		return disaster
	}
}

func (s *Simulation) generateActors() {
	if s.Disaster != 0 {
		return
	}

	numActors := 3 + rand.Intn(7)
	for i := 0; i < s.MakeSigns; i++ {
		numActors += 2 + rand.Intn(10)
	}
	for i := 0; i < numActors; i++ {
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

	s.NoSaleNoStock = 0

	s.Disaster = s.nextDisaster()
	s.DisasterAcknowledged = false

	s.generateActors()
}

func (s *Simulation) Tick() error {
	if s.DayFinished {
		return nil
	}

	standX, standY := 19, 10

	spaceAvailable := func(x, y int, thisActor *SimulationActor) bool {
		for _, a := range s.Actors {
			if a == thisActor || a.Inactive {
				continue
			} else if a.X == x && a.Y == y && (a.TargetActive == thisActor.TargetActive || (x == standX && y == standY)) {
				return false
			}
		}
		return true
	}

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
		moveX, moveY := 0, 0
		if a.X < tx {
			moveX++
		} else if a.X > tx {
			moveX--
		}
		if a.Y < ty {
			moveY++
		} else if a.Y > ty {
			moveY--
		}
		a.X, a.Y = a.X+moveX, a.Y+moveY

		if a.TargetActive && (a.X < -2 || a.X > 40 || a.Y > 17) {
			a.Inactive = true
			continue
		}

		if a.X == tx && a.Y == ty {
			if !a.TargetActive {
				if s.MakePretzels > s.PretzelsSold {
					s.PretzelsSold++
					a.WaitTicks = 100 + rand.Intn(150)
				} else {
					s.NoSaleNoStock++
					a.WaitTicks = 25 + rand.Intn(25)
				}

				a.TargetActive = true
				continue
			} else {
				a.Inactive = true
				continue
			}
		}
		a.WaitTicks = a.WalkTicks
		if !spaceAvailable(a.X+moveX, a.Y+moveY, a) || !spaceAvailable(a.X+moveX-1, a.Y+moveY, a) || !spaceAvailable(a.X+moveX+1, a.Y+moveY, a) {
			a.WaitTicks = 50 + rand.Intn(100)
		}
	}
	if active == 0 {
		s.DayFinishedTicks--
		if s.DayFinishedTicks == 0 {
			income := s.PretzelsSold * s.PretzelPrice
			expenses := (s.MakePretzels * PretzelCost) + (s.MakeSigns * SignCost)
			s.Money = s.Money + income - expenses

			s.DayFinished = true
		}
	}
	return nil
}
