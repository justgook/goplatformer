package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game/stage"
)

type Game struct {
	Stages []stage.Stage
	Stage  stage.Stage
	// UI           *ui.UI1
	CurrentWorld int
}

func New() *Game {
	// 640x360
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Game title goes here")

	current := &stage.Play{}
	current.Init()

	return &Game{
		Stage: current,
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	// 640x360
	return 640, 360 //- KEEP THAT ASPECT RATIO
	// return w >> 1, h >> 1
}
func (g *Game) Update() error {
	g.Stage.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 20, G: 20, B: 40, A: 255})
	g.Stage.Draw(screen)
}

