package game

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game/state"
)

type Game struct {
	State *state.GameState
}

func New() *Game {
	deviceInfo := &state.DeviceInfo{
		ScreenWidth:  640,
		ScreenHeight: 360,
	}

	gameState := &state.GameState{
		DeviceInfo: deviceInfo,
	}
	slog.Info("Update seed to random")
	gameState.Init(deviceInfo.ScreenWidth, deviceInfo.ScreenHeight, 1)

	gameState.SetScene(&BlackScene{})
	gameState.SetScene(&IntroScene{})

	gameState.SetScene(&PlayScene{})

	return &Game{
		State: gameState,
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.State.DeviceInfo.ScreenWidth, g.State.DeviceInfo.ScreenHeight
}
func (g *Game) Update() error {
	return g.State.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.State.Draw(screen)
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))

}
