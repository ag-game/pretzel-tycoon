package game

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"code.rocketnine.space/tslocum/etk"
	"code.rocketnine.space/tslocum/messeji"
	"code.rocketnine.space/tslocum/pretzel-tycoon/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Game struct {
	inputBuffer *etk.Input
	textBuffer  *dummyTextBuffer
}

var addedGame bool

func loadFont() font.Face {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	return face
}

type dummyTextBuffer struct {
	*etk.Text
}

func (t *dummyTextBuffer) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return false, nil
}

func (t *dummyTextBuffer) HandleKeyboard() (handled bool, err error) {
	return false, nil
}

func NewGame() (*Game, error) {
	etk.Style.TextColorLight = color.RGBA{255, 255, 255, 255}
	etk.Style.TextColorDark = color.RGBA{255, 255, 255, 255}
	etk.Style.InputBgColor = color.RGBA{0, 0, 0, 255}
	etk.Style.TextBgColor = color.RGBA{0, 0, 0, 255}
	etk.Style.TextFont = loadFont()

	g := &Game{
		inputBuffer: etk.NewInput("", "", func(text string) (handled bool) {
			log.Println("selected", text)
			return true
		}),
		textBuffer: &dummyTextBuffer{
			Text: etk.NewText("Hello world!"),
		},
	}

	// Configure text buffer.
	g.textBuffer.Field.SetHorizontal(messeji.AlignStart)
	g.textBuffer.Field.SetVertical(messeji.AlignStart)
	g.textBuffer.Field.SetPadding(0)

	// Create window input buffer and text buffer.
	w := etk.NewWindow()
	w.AddChild(g.inputBuffer)
	w.AddChild(g.textBuffer)

	etk.SetRoot(w)

	return g, nil
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	// Maintain constant internal resolution.
	if world.InternalScreenWidth != world.ScreenWidth || world.InternalScreenHeight != world.ScreenHeight {
		world.ScreenWidth, world.ScreenHeight = world.InternalScreenWidth, world.InternalScreenHeight
		etk.Layout(world.ScreenWidth, world.ScreenHeight)
		g.textBuffer.SetRect(image.Rect(0, 3, world.ScreenWidth, world.ScreenHeight))
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

	// Handle user input.
	err := etk.Update()
	if err != nil {
		log.Fatal(err)
	}

	// TODO only do this when the view buffer or input buffer changes

	pretzelsSold := 50
	pretzelPrice := "$.10"
	totalIncome := "$5.00"

	pretzelsMade := 50
	signsMade := 3
	totalExpenses := "$1.45"

	profit := "$3.55"
	assets := "$6.60"

	g.textBuffer.Clear()
	g.textBuffer.Write([]byte(fmt.Sprintf(viewFinancialReport, pretzelsSold, pretzelPrice, totalIncome, pretzelsMade, signsMade, totalExpenses, profit, assets)))

	g.textBuffer.Clear()
	for i := 0; i < 22; i++ {
		if i != 0 {
			g.textBuffer.Write([]byte("\n"))
		}
		g.textBuffer.Write([]byte(fmt.Sprintf("%d", i)))
	}
	// TODO fix trailing newline causing scroll bar to appear
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//screen.cl

	// Draw the text buffer over the hidden input buffer.
	g.textBuffer.Draw(screen)

	if world.Debug != 0 {
		debugText := fmt.Sprintf("TPS %0.0f\nFPS %0.0f",
			ebiten.ActualTPS(),
			ebiten.ActualFPS())
		ebitenutil.DebugPrintAt(screen, debugText, 2, 0)
	}
}

func (g *Game) Exit() {
	os.Exit(0)
}
