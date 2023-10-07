package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/justgook/goplatformer/pkg/game/state"
)

type DeviceInfo struct {
	ScreenWidth  int
	ScreenHeight int
}

type Game struct {
	State      *state.GameState
	DeviceInfo *DeviceInfo
}

func New() *Game {
	deviceInfo := &DeviceInfo{
		ScreenWidth:  640,
		ScreenHeight: 360,
	}

	gameState := &state.GameState{}
	gameState.Init(deviceInfo.ScreenWidth, deviceInfo.ScreenHeight)

	gameState.SetScene(&BlackScene{})
	gameState.SetScene(&ItroScene{})
	// gameState.SetScene(&StartScene{})
	// gameState.SetScene(&PlayScene{})

	return &Game{
		State:      gameState,
		DeviceInfo: deviceInfo,
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.DeviceInfo.ScreenWidth, g.DeviceInfo.ScreenHeight
}
func (g *Game) Update() error {
	g.State.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.State.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))

}

