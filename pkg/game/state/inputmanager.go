package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputState TODO update to https://github.com/quasilyte/ebitengine-input ???
type InputState struct {
	Down        bool
	JustPressed bool
}

type Input struct {
	N    *InputState
	E    *InputState
	S    *InputState
	W    *InputState
	Jump *InputState
}

type InputManager struct {
	*Input
}

func (i *InputManager) Init() {
	i.Input = &Input{}
	i.N = &InputState{}
	i.E = &InputState{}
	i.S = &InputState{}
	i.W = &InputState{}
	i.Jump = &InputState{}

}
func (i *InputManager) Update() {
	i.E.Down = ebiten.IsKeyPressed(ebiten.KeyRight) ||
		ebiten.GamepadAxisValue(0, 0) > 0.1

	i.W.Down = ebiten.IsKeyPressed(ebiten.KeyLeft) ||
		ebiten.GamepadAxisValue(0, 0) < -0.1

	i.S.Down = ebiten.IsKeyPressed(ebiten.KeyDown) ||
		ebiten.GamepadAxisValue(0, 1) > 0.1

	i.Jump.Down = ebiten.IsKeyPressed(ebiten.KeyX) ||
		ebiten.IsKeyPressed(ebiten.KeySpace) ||
		ebiten.IsGamepadButtonPressed(0, 0)
	i.Jump.JustPressed = inpututil.IsKeyJustPressed(ebiten.KeyX) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsGamepadButtonJustPressed(0, 0)
}
