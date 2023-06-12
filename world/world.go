package world

const (
	DefaultTPS = 60

	DefaultScreenWidth  = 1280
	DefaultScreenHeight = 720

	InternalScreenWidth  = 320
	InternalScreenHeight = 180

	MaxDebug = 2
)

var (
	TPS = DefaultTPS

	ScreenWidth  int
	ScreenHeight int

	Fullscreen   bool
	DisableVsync bool

	Debug int

	WASM bool
)
