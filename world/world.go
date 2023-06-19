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
	ViewStartDaySupplies
	ViewDay
	ViewFinancialReport
)

var (
	ScreenWidth  int
	ScreenHeight int

	StartingView ViewType
	Fullscreen   bool
	DisableVsync bool

	Debug int

	WASM bool
)
