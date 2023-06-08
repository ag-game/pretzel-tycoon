package world

const (
	DefaultTPS = 60

	DefaultScreenWidth  = 1280
	DefaultScreenHeight = 720

	InternalScreenWidth, InternalScreenHeight = 854, 480

	MaxDebug = 2
)

var (
	TPS = DefaultTPS

	ScreenWidth, ScreenHeight int

	Fullscreen   bool
	DisableVsync bool

	Debug int

	WASM bool
)
