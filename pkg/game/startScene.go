package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer"
	"github.com/justgook/goplatformer/pkg/game/state"
)

var _ state.Scene = (*StartScene)(nil)

type StartScene struct {
	tick int
}

// Draw implements Scene.
func (s *StartScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 5, G: 23, B: 38, A: 255})
	s.drawTitle(screen)
	content := "Space To Start"
	if (s.tick/20)%2 > 0 {
		content = "> " + content + " <"
	}
	goplatformer.DrawTextWithShadowCenter(
		screen,
		content,
		0,
		200,
		1,
		color.White,
		640,
	)
}

func (s *StartScene) drawTitle(screen *ebiten.Image) {
	img := goplatformer.StartMenuTitleImg
	screen.Fill(color.RGBA{R: 5, G: 23, B: 38, A: 255})
	op := &ebiten.DrawImageOptions{}
	x := screen.Bounds().Dx()/2 - img.Bounds().Dx()/2
	y := -64 + screen.Bounds().Dy()/2 - img.Bounds().Dy()/2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}

// Init implements Scene.
func (s *StartScene) Init() {
}

// Terminate implements Scene.
func (s *StartScene) Terminate() {
}

// Update implements Scene.
func (s *StartScene) Update(state *state.GameState) error {
	if state.Input.Jump.JustPressed {
		state.SetScene(&PlayScene{})
	}
	s.tick++
	return nil
}

