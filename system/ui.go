package system

import (
	"fmt"
	"strings"

	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/pretzel-tycoon/component"
	"code.rocketnine.space/tslocum/pretzel-tycoon/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type UISystem struct {
	*component.Once

	initialized bool

	tmpImg *ebiten.Image
}

func (u *UISystem) initialize() {
	u.tmpImg = ebiten.NewImage(world.InternalScreenWidth/2, world.InternalScreenHeight/2)
	u.initialized = true
}

func (u *UISystem) Update(e gohan.Entity) error {
	if !u.initialized {
		u.initialize()
	}
	return nil
}

func (u *UISystem) Draw(e gohan.Entity, screen *ebiten.Image) error {
	if !u.initialized {
		u.initialize()
	}

	if world.Debug != 0 {
		debugText := fmt.Sprintf("TPS %0.0f\nFPS %0.0f",
			ebiten.ActualTPS(),
			ebiten.ActualFPS())
		ebitenutil.DebugPrintAt(screen, debugText, 2, 0)
	}
	return nil
}

func (u *UISystem) drawCenteredText(screen *ebiten.Image, text string) {
	const (
		charWidth  = 6
		charHeight = 16
		textScale  = 2
	)

	lines := strings.Split(text, "\n")
	var w int
	for _, line := range lines {
		l := len(line)
		if l > w {
			w = l
		}
	}

	u.tmpImg.Clear()
	for i, line := range lines {
		x, y := ((w-len(line))/2)*charWidth, i*charHeight
		ebitenutil.DebugPrintAt(u.tmpImg, line, x, y)
	}

	width := float64(w) * charWidth * textScale
	height := float64(len(lines) * charHeight * textScale)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Scale(textScale, textScale)
	op.GeoM.Translate(world.InternalScreenWidth/2-width/2, world.InternalScreenHeight/2-height/2)
	screen.DrawImage(u.tmpImg, op)
}
