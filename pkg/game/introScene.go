package game

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer"
	"github.com/justgook/goplatformer/pkg/game/state"
)

var _ state.Scene = (*IntroScene)(nil)

type IntroScene struct {
	startTime time.Time
}

// Draw implements Scene.
func (*IntroScene) Draw(screen *ebiten.Image) {
	img := goplatformer.IntroLogo0x069
	screen.Fill(color.RGBA{R: 5, G: 23, B: 38, A: 255})
	op := &ebiten.DrawImageOptions{}
	x := screen.Bounds().Dx()/2 - img.Bounds().Dx()/2
	y := screen.Bounds().Dy()/2 - img.Bounds().Dy()/2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}

// Init implements Scene.
func (s *IntroScene) Init(st *state.GameState) {
	s.startTime = time.Now()
}

// Terminate implements Scene.
func (s *IntroScene) Terminate() {
}

// Update implements Scene.
func (s *IntroScene) Update(state *state.GameState) error {
	if time.Since(s.startTime).Seconds() > 2 {
		state.SetScene(&StartScene{})
	}
	return nil
}
