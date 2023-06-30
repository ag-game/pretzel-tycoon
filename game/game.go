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

	inputLetters bool // Whether to allow the user to input letters.

	sim *Simulation

	gameOver bool
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
		currentView: world.StartingView,
		dayBuffer:   make([][]byte, 18),
		sim:         NewSimulation(),
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
	for _, view := range world.InputViews {
		if g.currentView == view {
			return true
		}
	}
	return false
}

func (g *Game) acceptInput(text string) (handled bool) {
	if text == "" {
		switch g.currentView {
		case world.ViewStartDayProduction1:
			if g.sim.MakePretzels == -1 || g.sim.MakePretzels*PretzelCost > g.sim.Money {
				return false
			}
		case world.ViewStartDayProduction2:
			if g.sim.MakeSigns == -1 || g.sim.MakeSigns*SignCost > (g.sim.Money-(g.sim.MakePretzels*PretzelCost)) {
				return false
			}
		case world.ViewStartDayProduction3:
			if g.sim.PretzelPrice == -1 {
				return false
			}
		}
	} else {
		// TODO handle non-numeric input
		i, err := strconv.Atoi(text)
		if err != nil {
			return false
		}
		switch g.currentView {
		case world.ViewStartDayProduction1:
			if i < 1 || i*PretzelCost > g.sim.Money {
				return false
			}
			g.sim.MakePretzels = i
		case world.ViewStartDayProduction2:
			if i*SignCost > (g.sim.Money - (g.sim.MakePretzels * PretzelCost)) {
				return false
			}
			g.sim.MakeSigns = i
		case world.ViewStartDayProduction3:
			g.sim.PretzelPrice = i
		}
	}

	g.currentView++
	partialTransition := g.currentView == world.ViewStartDayProduction2 || g.currentView == world.ViewStartDayProduction3
	if partialTransition {
		viewBytes := viewText[g.currentView-1]
		lines := bytes.Split(viewBytes, []byte("\n"))
		g.viewTicks = len(lines) - 2
	} else {
		g.viewTicks = 0
	}

	if g.currentView == world.ViewDay {
		g.sim.StartDay()
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
	for _, a := range g.sim.Actors {
		g.setDayCell(a.X, a.Y, 'o')
		g.setDayCell(a.X+1, a.Y, 'o')
	}

	// Draw labels.
	drawText(17, 0, fmt.Sprintf("DAY %d", g.sim.Day))
	drawText(8, 2, fmt.Sprintf("STOCK %d", g.sim.MakePretzels-g.sim.PretzelsSold))
	drawText(24, 2, fmt.Sprintf("SOLD %d", g.sim.PretzelsSold))
	if g.sim.Disaster != 0 {
		label := disasterLabel(g.sim.Disaster)
		drawText(20-len(label)/2, 12, label)

		drawText(11, 14, "NO CUSTOMERS TODAY")
		drawText(7, 17, "PRESS ENTER TO CONTINUE...")
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

	if g.gameOver {
		g.textBuffer.Clear()
		g.textBuffer.Write([]byte(fmt.Sprintf(`                 DAY %d




           YOU HAVE NO MONEY!

          YOU ARE NOW HOMELESS!

          THERE IS NO HOPE LEFT
       IN YOUR EYES... ONLY PAIN!






                GAME OVER`, g.sim.Day)))
		return nil
	}

	if g.currentView == world.ViewDay {
		return g.drawDay()
	}

	income := g.sim.PretzelsSold * g.sim.PretzelPrice
	totalIncome := fmt.Sprintf("$%d.%02d", income/100, income%100)

	expenses := (g.sim.MakePretzels * PretzelCost) + (g.sim.MakeSigns * SignCost)
	totalExpenses := formatMoney(expenses)

	profit := formatMoney(income - expenses)
	assets := formatMoney(g.sim.Money)

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

	formatCents := func(cents int) string {
		return fmt.Sprintf("$%d.%02d", cents/100, cents%100)
	}

	// Format view.
	var lines [][]byte
	switch g.currentView {
	case world.ViewStartDayProduction1:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.sim.Day, formatMoney(g.sim.Money)))
		lines = bytes.Split(viewBytes, []byte("\n"))
	case world.ViewStartDayProduction2:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.sim.Day, formatMoney(g.sim.Money), g.sim.MakePretzels))
		lines = bytes.Split(viewBytes, []byte("\n"))
	case world.ViewStartDayProduction3:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.sim.Day, formatMoney(g.sim.Money), g.sim.MakePretzels, g.sim.MakeSigns))
		lines = bytes.Split(viewBytes, []byte("\n"))
	case world.ViewFinancialReport:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.sim.Day, g.sim.PretzelsSold, formatCents(g.sim.PretzelPrice), totalIncome, g.sim.MakePretzels, g.sim.MakeSigns, totalExpenses, profit, assets))
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

	if g.gameOver {
		return nil
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
			g.sim.Day++
			g.currentView = world.ViewStartDayProduction1
		} else {
			// Skip to end of day.
			if g.currentView == world.ViewDay {
				for !g.sim.DayFinished {
					err := g.sim.Tick()
					if err != nil {
						return err
					}
				}
				g.sim.DisasterAcknowledged = true

				if g.sim.Money < PretzelCost {
					g.gameOver = true
				}
			}
			g.currentView++
		}
		g.viewTicks = 0

		err := g.refreshBuffer()
		if err != nil {
			return err
		}
		return nil
	}

	if g.currentView == world.ViewDay {
		numTicks := 1
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			numTicks = 4
		}
		for i := 0; i < numTicks; i++ {
			err := g.sim.Tick()
			if err != nil {
				return err
			}
		}

		if g.sim.DayFinished && (g.sim.Disaster == 0 || g.sim.DisasterAcknowledged) {
			g.currentView++
			g.viewTicks = 0

			if g.sim.Money < PretzelCost {
				g.gameOver = true
			}
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

func formatMoney(v int) string {
	var prefix string
	if v < 0 {
		prefix = "-"
		v *= -1
	}
	return fmt.Sprintf("%s$%d.%02d", prefix, v/100, v%100)
}
