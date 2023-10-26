package state

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const transitionMaxCount = 20

type Scene interface {
	Init(state *GameState)
	Update(state *GameState) error
	Draw(screen *ebiten.Image)
	Terminate()
}

type SceneManager struct {
	current         Scene
	next            Scene
	transitionCount int
	transitionFrom  *ebiten.Image
	transitionTo    *ebiten.Image
	runner          struct {
		Update func(state *GameState) error
		Draw   func(r *ebiten.Image)
	}
}
