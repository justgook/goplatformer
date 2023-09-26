package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game/stage"
	"image/color"
)

type Game struct {
	Stages        []stage.Stage
	Stage         stage.Stage
	CurrentWorld  int
	Width, Height int
}

func New() *Game {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Game title goes here")

	current := &stage.Play{}
	current.Init()

	return &Game{
		Stage: current,
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}
func (g *Game) Update() error {
	g.Stage.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 20, G: 20, B: 40, A: 255})
	g.Stage.Draw(screen)
}

