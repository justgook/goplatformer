package game

import (
	"fmt"
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

	ui *ui.View
}

// Init implements state.Scene.
func (p *PlayScene) Init(st *state.GameState) {
	p.world = ecs.NewECS(donburi.NewWorld())

	collisionSpace := components.Space.Get(entity.CreateSpace(p.world))

	levelEntity := entity.CreateLevel(p.world, st, collisionSpace)
	level := components.Level.Get(levelEntity)

	// Spawn player START
	playerEnt := entity.CreatePlayer(p.world, *goplatformer.EmbeddedPlayerAnimation, st.Input.Input)
	player := components.Object.Get(playerEnt)
	player.X = float64(128)
	player.Y = float64(128)
	collisionSpace.Add(player)
	// Spawn player END

	// TEST BULLET EMITTER//
	playerEntry := playerEnt.Entity()
	entity.CreateBulletEmitter(p.world, goplatformer.TestBullet, playerEntry, playerEntry)

	// CAMERA STUFF Start

	p.worldScreen = ebiten.NewImage(level.Size.X, level.Size.Y)
	p.camera = &components.CameraData{ViewPort: f64.Vec2{float64(st.DeviceInfo.ScreenWidth), float64(st.DeviceInfo.ScreenHeight)}}
	p.camera.SetPosition(f64.Vec2{})

	// CAMERA STUFF End

	//ui.Debug = true

	// -------------------- UI STUFF START ----------------------------
	p.ui = &ui.View{
		Width:  st.DeviceInfo.ScreenWidth,
		Height: st.DeviceInfo.ScreenHeight,
	}

	levelMap := &panel.LevelMap{}

	levelMap.Init(components.Level.Get(levelEntity), &level.CurrentRoomXY)
	p.ui.AddChild(levelMap.View)
	// ---------------------- UI STUFF END --------------------------

	/*-----------------------------------SYSTEMS-------------------------------------------*/
	p.world.AddSystem(system.UpdatePlayer)
	p.world.AddSystem(system.UpdateCharAnim)
	p.world.AddSystem(system.UpdateBulletEmitter)
	p.world.AddSystem(system.UpdateBullets)

	p.world.AddSystem(system.EnemyTrigger)
	p.world.AddSystem(system.ExitTrigger)

	/*-----------------------------------RENDERS-------------------------------------------*/
	p.world.AddRenderer(layers.Background, renderer.Level)
	p.world.AddRenderer(layers.Background, renderer.CharAnim)
	p.world.AddRenderer(layers.Background, renderer.Bullets)

	//p.world.AddRenderer(layers.Background, renderer.Wall)
}

// Draw implements state.Scene.
func (p *PlayScene) Draw(screen *ebiten.Image) {
	p.worldScreen.Clear()
	p.world.Draw(p.worldScreen)
	p.camera.Render(p.worldScreen, screen)

	p.ui.Draw(screen)
	//DEBUG INFO
	p.camera.DebugInfoDraw(screen)

	levelEnt, _ := components.Level.First(p.world.World)
	levelComp := components.Level.Get(levelEnt)
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Room: %d,%d",
			levelComp.CurrentRoomXY.X, levelComp.CurrentRoomXY.Y,
		),
		int(p.camera.ViewPort[0])-96, int(p.camera.ViewPort[1])-16,
	)
}

// Update implements state.Scene.
func (p *PlayScene) Update(state *state.GameState) error {
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
