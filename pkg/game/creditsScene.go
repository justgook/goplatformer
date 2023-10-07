package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game/state"
)

var _ state.Scene = (*CreditsScene)(nil)

type CreditsScene struct {
}

// Draw implements Scene.
func (*CreditsScene) Draw(screen *ebiten.Image) {
}

// Init implements Scene.
func (*CreditsScene) Init() {
}

// Terminate implements Scene.
func (*CreditsScene) Terminate() {
}

// Update implements Scene.
func (*CreditsScene) Update(state *state.GameState) error {
	return nil
}

