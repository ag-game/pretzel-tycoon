package game

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

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

var matchNumbers = regexp.MustCompile("^[0-9]+$")

type Game struct {
	inputBuffer *etk.Input
	textBuffer  *etk.Text

	currentView world.ViewType
	viewTicks   int

	dayBuffer [][]byte

	day int

	inputLetters bool // Whether to allow the user to input letters.

	makePretzels     int // In dozens.
	makePretzelsLast int
	makeSigns        int
	makeSignsLast    int

	pretzelPrice     int // In cents.
	pretzelPriceLast int

	simulation *Simulation
}

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

func NewGame() (*Game, error) {
	rand.Seed(time.Now().UnixNano())

	etk.Style.TextColorLight = color.RGBA{255, 255, 255, 255}
	etk.Style.TextColorDark = color.RGBA{255, 255, 255, 255}
	etk.Style.InputBgColor = color.RGBA{0, 0, 0, 255}
	etk.Style.InputBgColor = color.RGBA{0, 0, 0, 255}
	etk.Style.TextBgColor = color.RGBA{0, 0, 0, 255}
	etk.Style.TextFont = loadFont()

	g := &Game{
		currentView:      world.StartingView,
		day:              1,
		makePretzels:     -1,
		makePretzelsLast: -1,
		makeSigns:        -1,
		makeSignsLast:    -1,
		pretzelPrice:     -1,
		pretzelPriceLast: -1,
		dayBuffer:        make([][]byte, 18),
		simulation:       &Simulation{},
	}
	g.inputBuffer = etk.NewInput("", "", g.acceptInput)
	g.textBuffer = etk.NewText("Hello world!")

	for i := range g.dayBuffer {
		g.dayBuffer[i] = make([]byte, 40)
	}

	// Configure text buffer.
	g.textBuffer.Field.SetHorizontal(messeji.AlignStart)
	g.textBuffer.Field.SetVertical(messeji.AlignStart)
	g.textBuffer.Field.SetPadding(0)
	g.textBuffer.Field.SetLineHeight(10)

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
		g.textBuffer.SetRect(image.Rect(0, 1, world.ScreenWidth, world.ScreenHeight))
	}
	return world.ScreenWidth, world.ScreenHeight
}

func (g *Game) inputActive() bool {
	switch g.currentView {
	case world.ViewStartDayProduction1:
		return g.makePretzels == -1
	case world.ViewStartDayProduction2:
		return g.makeSigns == -1
	case world.ViewStartDayProduction3:
		return g.pretzelPrice == -1
	}
	return false
}

func (g *Game) acceptInput(text string) (handled bool) {
	if text == "" {
		switch g.currentView {
		case world.ViewStartDayProduction1:
			if g.makePretzelsLast == -1 {
				return false
			}
			g.makePretzels = g.makePretzelsLast
		case world.ViewStartDayProduction2:
			if g.makeSignsLast == -1 {
				return false
			}
			g.makeSigns = g.makeSignsLast
		case world.ViewStartDayProduction3:
			if g.pretzelPriceLast == -1 {
				return false
			}
			g.pretzelPrice = g.pretzelPriceLast
		}
	} else {
		// TODO handle non-numeric input
		i, err := strconv.Atoi(text)
		if err != nil {
			return false
		}
		switch g.currentView {
		case world.ViewStartDayProduction1:
			g.makePretzels = i
		case world.ViewStartDayProduction2:
			g.makeSigns = i
		case world.ViewStartDayProduction3:
			g.pretzelPrice = i
		}
	}

	g.currentView++
	partialTransition := g.currentView == world.ViewStartDayProduction2 || g.currentView == world.ViewStartDayProduction3
	if partialTransition {
		viewBytes := viewText[g.currentView-1]
		lines := bytes.Split(viewBytes, []byte("\n"))
		g.viewTicks = len(lines) - 4
	} else {
		g.viewTicks = 0
	}

	log.Println(g.currentView)
	if g.currentView == world.ViewDay {
		g.simulation.StartDay()
	}

	g.refreshBuffer()
	return true
}

func (g *Game) setDayCell(x int, y int, c byte) error {
	if x < 0 || y < 0 || x > 39 || y > 17 {
		// Skip drawing off-screen characters.
		return nil
	}
	g.dayBuffer[y][x] = c
	return nil
}

func (g *Game) drawDay() error {
	for y := range g.dayBuffer {
		for x := range g.dayBuffer[y] {
			g.dayBuffer[y][x] = ' '
		}
	}

	drawLine := func(x int, y int, width int) {
		for i := 0; i < width; i++ {
			g.setDayCell(x+i, y, '_')
		}
	}

	drawText := func(x int, y int, text string) {
		for i, r := range text {
			g.setDayCell(x+i, y, byte(r))
		}
	}

	// Draw pretzel stand.
	drawPretzelStand := func(x int, y int) {
		width := 14
		height := 3

		// Draw outline.
		drawLine(x+1, y, width)
		for cy := 1; cy <= height; cy++ {
			g.setDayCell(x, y+cy, '|')
			g.setDayCell(x+width+1, y+cy, '|')
		}
		drawLine(x+1, y+height, width)

		// Draw sign text.
		drawText(x+2, y+2, "& PRETZELS &")
	}
	drawPretzelStand(12, 6)

	// Draw actors.
	for _, a := range g.simulation.Actors {
		g.setDayCell(a.X, a.Y, 'o')
		g.setDayCell(a.X+1, a.Y, 'o')
	}

	g.textBuffer.Clear()
	for y := range g.dayBuffer {
		if y != 0 {
			g.textBuffer.Write([]byte("\n"))
		}
		g.textBuffer.Write(g.dayBuffer[y])
	}
	return nil
}

func (g *Game) refreshBuffer() error {
	// TODO only do this when the view buffer or input buffer changes
	// TODO fix trailing newline causing scroll bar to appear

	if g.currentView == world.ViewDay {
		return g.drawDay()
	}

	pretzelsSold := 50
	pretzelPrice := "$.10"
	totalIncome := "$5.00"

	pretzelsMade := 50
	signsMade := 3
	totalExpenses := "$1.45"

	profit := "$3.55"
	assets := "$6.60"

	viewBytes := viewText[g.currentView]
	writeLines := g.viewTicks + 1

	// Append start screen text.
	if g.currentView == world.ViewTitle && g.viewTicks%200 < 150 {
		viewBytes = append(viewBytes, bytes.TrimRight(centeredText("PRESS ENTER TO START"), "\n")...)
	}

	// Append user input.
	if g.inputActive() {
		viewBytes = append(viewBytes, g.inputBuffer.Text()...)

		// Append cursor icon.
		if g.viewTicks%150 < 100 {
			viewBytes = append(viewBytes, '|')
		}
	}

	// Format view.
	var lines [][]byte
	switch g.currentView {
	case world.ViewStartDayProduction1:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.day))
		lines = bytes.Split(viewBytes, []byte("\n"))
	case world.ViewStartDayProduction2:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.day, g.makePretzels))
		lines = bytes.Split(viewBytes, []byte("\n"))
	case world.ViewStartDayProduction3:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.day, g.makePretzels, g.makeSigns))
		lines = bytes.Split(viewBytes, []byte("\n"))
	case world.ViewFinancialReport:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.day, pretzelsSold, pretzelPrice, totalIncome, pretzelsMade, signsMade, totalExpenses, profit, assets))
		lines = bytes.Split(viewBytes, []byte("\n"))
	default:
		lines = bytes.Split(viewBytes, []byte("\n"))
	}

	g.textBuffer.Clear()
	if writeLines <= 0 {
		return nil
	} else if writeLines >= len(lines) {
		g.textBuffer.Write(viewBytes)
		return nil
	}
	wrote := 0
	for i := 0; i < len(lines); i++ {
		if i != 0 {
			g.textBuffer.Write([]byte("\n"))
		}
		if len(lines[i]) == 0 {
			continue
		} else if wrote == writeLines {
			return nil
		}
		g.textBuffer.Write(lines[i])
		wrote++
	}
	return nil
}

func (g *Game) resetDay() error {
	g.makePretzelsLast = g.makePretzels
	g.makePretzels = -1

	g.makeSignsLast = g.makeSigns
	g.makeSigns = -1

	g.pretzelPriceLast = g.pretzelPrice
	g.pretzelPrice = -1
	return nil
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
	if g.inputActive() {
		err := etk.Update()
		if err != nil {
			log.Fatal(err)
		}

		inputText := g.inputBuffer.Text()
		if len(inputText) > 0 && !g.inputLetters && !matchNumbers.MatchString(inputText) {
			var newInput string
			for _, r := range inputText {
				if matchNumbers.MatchString(string(r)) {
					newInput += string(r)
				}
			}
			g.inputBuffer.Clear()
			g.inputBuffer.Write([]byte(newInput))
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyKPEnter) || (g.currentView != world.ViewDay && inpututil.IsKeyJustPressed(ebiten.KeySpace)) {
		if g.currentView == world.ViewFinancialReport {
			g.resetDay()
			g.currentView = world.ViewStartDayProduction1
		} else {
			g.currentView++
		}
		g.viewTicks = 0
	}

	if g.currentView == world.ViewDay {
		numTicks := 1
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			numTicks = 4
		}
		for i := 0; i < numTicks; i++ {
			err := g.simulation.Tick()
			if err != nil {
				return err
			}
		}

		if g.simulation.DayFinished {
			g.currentView++
			g.viewTicks = 0
		}
	}

	err := g.refreshBuffer()
	if err != nil {
		return err
	}

	g.viewTicks++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
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
