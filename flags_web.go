//go:build js && wasm
// +build js,wasm

package main

import (
	"code.rocketnine.space/tslocum/pretzel-tycoon/world"
	"github.com/hajimehoshi/ebiten/v2"
)

func parseFlags() {
	world.WASM = true

	ebiten.SetFullscreen(true)
}
