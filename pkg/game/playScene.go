package game

import (
	"fmt"
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/justgook/goplatformer"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/components/panel"
	"github.com/justgook/goplatformer/pkg/game/entity"
	"github.com/justgook/goplatformer/pkg/game/layers"
	"github.com/justgook/goplatformer/pkg/game/renderer"
	"github.com/justgook/goplatformer/pkg/game/state"
	"github.com/justgook/goplatformer/pkg/game/system"
	"github.com/justgook/goplatformer/pkg/ui"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"golang.org/x/image/math/f64"
)

var _ state.Scene = (*PlayScene)(nil)

type PlayScene struct {
	world *ecs.ECS

	worldScreen *ebiten.Image
	camera      *components.CameraData

	delme image.Point

	ui *ui.View
}

// Init implements state.Scene.
func (p *PlayScene) Init(st *state.GameState) {
	p.world = ecs.NewECS(donburi.NewWorld())

	spaceEntity := entity.CreateSpace(p.world)
	space := components.Space.Get(spaceEntity)

	levelEntity := entity.CreateLevel(p.world, st, space)
	level := components.Level.Get(levelEntity)

	// Spawn player START
	playerEnt := entity.CreatePlayer(p.world, *goplatformer.EmbeddedPlayerAnimation, st.Input.Input)
	player := components.Object.Get(playerEnt)
	player.X = float64(level.Start.X*256 + 128)
	player.Y = float64(level.Start.Y*256 + 128)
	space.Add(player)
	// Spawn player END

	//entity.CreatePlayer(p.world, *goplatformer.EmbeddedPlayerAnimation, st.Input.Input),

	// CAMERA STUFF Start

	p.worldScreen = ebiten.NewImage(level.Size.X, level.Size.Y)
	p.camera = &components.CameraData{
		ViewPort: f64.Vec2{
			float64(st.DeviceInfo.ScreenWidth),
			float64(st.DeviceInfo.ScreenHeight),
		},
	}
	p.camera.SetPosition(
		f64.Vec2{
			float64(components.Level.Get(levelEntity).Start.X*256 + 128),
			float64(components.Level.Get(levelEntity).Start.Y*256 + 128),
		},
	)
	slog.Info("START",
		"X", components.Level.Get(levelEntity).Start.X,
		"Y", components.Level.Get(levelEntity).Start.Y,
	)
	// CAMERA STUFF End

	//ui.Debug = true

	// -------------------- UI STUFF START ----------------------------
	p.ui = &ui.View{
		Width:  st.DeviceInfo.ScreenWidth,
		Height: st.DeviceInfo.ScreenHeight,
	}

	levelMap := &panel.LevelMap{}
	levelMap.Init(components.Level.Get(levelEntity), &p.delme)
	p.ui.AddChild(levelMap.View)
	// ---------------------- UI STUFF END --------------------------

	p.world.AddSystem(system.UpdatePlayer)
	p.world.AddSystem(system.UpdateCharAnim)
	p.world.AddRenderer(layers.Background, renderer.Level)
	//p.world.AddRenderer(layers.Background, renderer.Player)
	p.world.AddRenderer(layers.Background, renderer.CharAnim)

	//p.world.AddRenderer(layers.Background, renderer.Wall)
}

// Draw implements state.Scene.
func (p *PlayScene) Draw(screen *ebiten.Image) {
	p.world.Draw(p.worldScreen)
	p.camera.Render(p.worldScreen, screen)

	p.ui.Draw(screen)
	//DEBUG INFO
	p.camera.DebugInfoDraw(screen)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Room: %d,%d",
			p.delme.X, p.delme.Y,
		),
		int(p.camera.ViewPort[0])-96, int(p.camera.ViewPort[1])-16,
	)
}

// Update implements state.Scene.
func (p *PlayScene) Update(state *state.GameState) error {
	p.delme = p.cameraToRoom()
	p.camera.DebugUpdate()
	p.cameraFollowPlayer()

	p.world.Update()
	p.ui.Update()

	return nil
}

func (p *PlayScene) cameraFollowPlayer() {
	playerEnt, _ := components.Player.First(p.world.World)
	playerObj := components.Object.Get(playerEnt)

	p.camera.SetPosition(f64.Vec2{playerObj.X, playerObj.Y})
}

// Terminate implements state.Scene.
func (p *PlayScene) Terminate() {
	p.world = nil
}

func (p *PlayScene) cameraToRoom() image.Point {
	output := p.camera.Center()
	return image.Point{
		X: int(output[0]) / 256, Y: int(output[1]) / 256,
	}
}
