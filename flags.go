//go:build !js || !wasm
// +build !js !wasm

package main

import (
	"flag"

	"code.rocketnine.space/tslocum/pretzel-tycoon/world"
)

func parseFlags() {
	var startingView int
	flag.BoolVar(&world.Fullscreen, "fullscreen", false, "run in fullscreen mode")
	flag.BoolVar(&world.DisableVsync, "no-vsync", false, "do not enable vsync (allows the game to run at maximum fps)")
	flag.IntVar(&world.Debug, "debug", 0, "debug level (0 - disabled, 1 - print fps and net stats, 2 - draw hitboxes)")
	flag.IntVar(&startingView, "view", 0, "start at specific view screen")

	flag.Parse()

	world.StartingView = world.ViewType(startingView)
}
