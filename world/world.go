package world

const (
	DefaultTPS = 100

	DefaultScreenWidth  = 1280
	DefaultScreenHeight = 720

	InternalScreenWidth  = 320
	InternalScreenHeight = 180

	MaxDebug = 2
)

type ViewType int

const (
	ViewTitle = iota
	ViewIntro1
	ViewStartDayProduction1
	ViewStartDayProduction2
	ViewStartDayProduction3
	ViewDay
	ViewFinancialReport
)

var InputViews = []ViewType{
	ViewStartDayProduction1,
	ViewStartDayProduction2,
	ViewStartDayProduction3,
}

var (
	ScreenWidth  int
	ScreenHeight int

	StartingView ViewType
	Fullscreen   bool
	DisableVsync bool

	Debug int

	WASM bool
)
