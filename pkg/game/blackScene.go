package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game/state"
)

var _ state.Scene = (*BlackScene)(nil)

type BlackScene struct {
}

// Draw implements Scene.
func (*BlackScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 5, G: 23, B: 38, A: 255})
}

// Init implements Scene.
func (*BlackScene) Init(st *state.GameState) {
}

// Terminate implements Scene.
func (*BlackScene) Terminate() {
}

// Update implements Scene.
func (*BlackScene) Update(state *state.GameState) error {
	return nil
}
