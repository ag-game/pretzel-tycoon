package game

import (
	"log"
	"os"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/pretzel-tycoon/entity"
	"code.rocketnine.space/tslocum/pretzel-tycoon/system"
	"code.rocketnine.space/tslocum/pretzel-tycoon/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
}

var addedGame bool

func NewGame() (*Game, error) {
	g := &Game{}

	if !addedGame {
		// Set up entity component system.
		entity.NewOnceEntity()
		gohan.AddSystem(&system.UISystem{})

		addedGame = true
	}

	return g, nil
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	// Maintain constant internal resolution.
	if world.InternalScreenWidth != world.ScreenWidth || world.InternalScreenHeight != world.ScreenHeight {
		world.ScreenWidth, world.ScreenHeight = world.InternalScreenWidth, world.InternalScreenHeight
	}
	return world.ScreenWidth, world.ScreenHeight
}

func (g *Game) Update() error {
	if ebiten.IsWindowBeingClosed() || (!world.WASM && ebiten.IsKeyPressed(ebiten.KeyEscape)) {
		g.Exit()
		return nil
	}

	// Toggle fullscreen.
	if (inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyKPEnter)) && ebiten.IsKeyPressed(ebiten.KeyAlt) {
		world.Fullscreen = !world.Fullscreen
		ebiten.SetFullscreen(world.Fullscreen)
		return nil
	}

	// Change debug level.
	if inpututil.IsKeyJustPressed(ebiten.KeyV) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		world.Debug++
		if world.Debug > world.MaxDebug {
			world.Debug = 0
		}
	}

	err := gohan.Update()
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	err := gohan.Draw(screen)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Exit() {
	os.Exit(0)
}
