package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game/stage"
)

type Game struct {
	Stages        []stage.Stage
	Stage         stage.Stage
	// UI           *ui.UI1
	CurrentWorld  int
}

func New() *Game {
	ebiten.SetWindowSize(512, 512)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Game title goes here")

	current := &stage.Play{}
	current.Init()

	// currentUI := &ui.UI1{}
	// currentUI.Init()
	return &Game{
		Stage: current,
		// UI: currentUI,
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}
func (g *Game) Update() error {
	g.Stage.Update()
	// g.UI.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 20, G: 20, B: 40, A: 255})
	g.Stage.Draw(screen)
	// g.UI.Draw(screen)
}
