/*
based on https://github.com/SolarLune/resolv.git
*/

package main

import (
	_ "embed"
	"errors"
	"fmt"
	"image/color"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/justgook/goplatformer"
	"github.com/justgook/goplatformer/pkg/resolv"
	"golang.org/x/image/font"
)

type Game struct {
	Worlds        []StageInterface
	Stage         *StageInterface
	CurrentWorld  int
	Width, Height int
	Debug         bool
	ShowHelpText  bool
	FontFace      font.Face
}

func NewGame() *Game {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("go Platformer")

	g := &Game{
		Width:        640,
		Height:       360,
		ShowHelpText: false,
		Debug:        false,
	}

	g.Worlds = []StageInterface{
		NewWorldPlatformer(g),
		NewWorldBouncer(g),
		NewWorldLineTest(g),
		NewWorldShapeTest(g),
		NewWorldDirectTest(g),
	}

	fontData, _ := truetype.Parse(goplatformer.ExcelFont)
	g.FontFace = truetype.NewFace(fontData, &truetype.Options{Size: 10})

	// Debug FPS rendering
	//go func() {
	//	for {
	//		fmt.Println("FPS: ", ebiten.ActualFPS())
	//		fmt.Println("Ticks: ", ebiten.ActualTPS())
	//		time.Sleep(time.Second)
	//	}
	//}()

	return g
}

func (g *Game) Update() error {
	var quit error

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if ebiten.ActualTPS() >= 60 {
			ebiten.SetTPS(6)
		} else {
			ebiten.SetTPS(60)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.Debug = !g.Debug
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		g.ShowHelpText = !g.ShowHelpText
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.CurrentWorld++
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.CurrentWorld--
	}

	if g.CurrentWorld >= len(g.Worlds) {
		g.CurrentWorld = 0
	} else if g.CurrentWorld < 0 {
		g.CurrentWorld = len(g.Worlds) - 1
	}

	world := g.Worlds[g.CurrentWorld]

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		world.Init()
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		quit = errors.New("quit")
	}

	world.Update()

	return quit
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 20, G: 20, B: 40, A: 255})
	g.Worlds[g.CurrentWorld].Draw(screen)
}

func (g *Game) DrawText(screen *ebiten.Image, x, y int, textLines ...string) {
	rectHeight := 10
	for _, txt := range textLines {
		w := float32(font.MeasureString(g.FontFace, txt).Round())

		vector.DrawFilledRect(
			screen,
			float32(x),
			float32(y-8),
			w,
			float32(rectHeight),
			color.RGBA{A: 192},
			false,
		)

		text.Draw(screen, txt, g.FontFace, x+1, y+1, color.RGBA{B: 150, A: 255})
		text.Draw(screen, txt, g.FontFace, x, y, color.RGBA{R: 100, G: 150, B: 255, A: 255})
		y += rectHeight
	}
}

func (g *Game) DebugDraw(screen *ebiten.Image, space *resolv.Space) {
	for y := 0; y < space.Height(); y++ {
		for x := 0; x < space.Width(); x++ {
			cell := space.Cell(x, y)

			cw := float32(space.CellWidth)
			ch := float32(space.CellHeight)

			cx := float32(cell.X) * cw
			cy := float32(cell.Y) * ch

			if cell.Occupied() {
				drawColor := color.RGBA{R: 255, G: 255, A: 255}
				vector.StrokeLine(screen, cx, cy, cx+cw, cy, 1, drawColor, false)
				vector.StrokeLine(screen, cx+cw, cy, cx+cw, cy+ch, 1, drawColor, false)
				vector.StrokeLine(screen, cx+cw, cy+ch, cx, cy+ch, 1, drawColor, false)
				vector.StrokeLine(screen, cx, cy, cx, cy+ch, 1, drawColor, false)
			}

		}
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.Width, g.Height
}

func main() {
	if err := ebiten.RunGame(NewGame()); err != nil {
		fmt.Print(err)
		return
	}
}
