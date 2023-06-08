//go:build js && wasm
// +build js,wasm

package main

import (
	"code.rocketnine.space/tslocum/pretzel-tycoon/world"
)

func parseFlags() {
	world.WASM = true

	ebiten.SetFullscreen(true)
}
