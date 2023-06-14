package game

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"regexp"
	"strconv"

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

type viewType int

const (
	viewTitle = iota
	viewIntro1
	viewStartDayProduction1
	viewStartDayProduction2
	viewStartDayProduction3
	viewStartDaySupplies
	viewFinancialReport
)

var matchNumbers = regexp.MustCompile("^[0-9]+$")

type Game struct {
	inputBuffer *etk.Input
	textBuffer  *dummyTextBuffer

	currentView viewType
	viewTicks   int

	day int

	inputLetters bool // Whether to allow the user to input letters.

	makePretzels int // In dozens.
	makeSigns    int

	pretzelPrice int // In cents.
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
		day: 1,
	}
	g.inputBuffer = etk.NewInput("", "", g.acceptInput)
	g.textBuffer = &dummyTextBuffer{
		Text: etk.NewText("Hello world!"),
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
	return g.currentView == viewStartDayProduction1 || g.currentView == viewStartDayProduction2 || g.currentView == viewStartDayProduction3
}

func (g *Game) acceptInput(text string) (handled bool) {
	log.Println("selected", text)

	i, err := strconv.Atoi(text)
	if err != nil {
		log.Println(err) // TODO this shouldnt happen
		return false
	}

	switch g.currentView {
	case viewStartDayProduction1:
		g.makePretzels = i
	case viewStartDayProduction2:
		g.makeSigns = i
	case viewStartDayProduction3:
		g.pretzelPrice = i
	}

	g.currentView++
	partialTransition := g.currentView == viewStartDayProduction2 || g.currentView == viewStartDayProduction3
	if partialTransition {
		viewBytes := viewText[g.currentView-1]
		lines := bytes.Split(viewBytes, []byte("\n"))
		g.viewTicks = len(lines) - 2
	} else {
		g.viewTicks = 0
	}
	return true
}

func (g *Game) refreshBuffer() error {
	// TODO only do this when the view buffer or input buffer changes
	// TODO fix trailing newline causing scroll bar to appear

	pretzelsSold := 50
	pretzelPrice := "$.10"
	totalIncome := "$5.00"

	pretzelsMade := 50
	signsMade := 3
	totalExpenses := "$1.45"

	profit := "$3.55"
	assets := "$6.60"

	viewBytes := viewText[g.currentView]
	writeLines := g.viewTicks - 1

	// Append start screen text.
	if g.currentView == viewTitle && g.viewTicks%200 < 150 {
		viewBytes = append(viewBytes, bytes.TrimRight(centeredText("PRESS SPACE TO START"), "\n")...)
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
	case viewStartDayProduction1:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.day))
		lines = bytes.Split(viewBytes, []byte("\n"))
	case viewStartDayProduction2:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.day, g.makePretzels))
		lines = bytes.Split(viewBytes, []byte("\n"))
	case viewStartDayProduction3:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.day, g.makePretzels, g.makeSigns))
		lines = bytes.Split(viewBytes, []byte("\n"))
	case viewStartDaySupplies:
		viewBytes = []byte(fmt.Sprintf(string(viewBytes), g.day))
		lines = bytes.Split(viewBytes, []byte("\n"))
	case viewFinancialReport:
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
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.currentView++
		if g.currentView > viewFinancialReport {
			g.currentView = viewTitle
		}
		g.viewTicks = 0
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
