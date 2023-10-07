package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer"
	"github.com/justgook/goplatformer/pkg/game/state"
	"github.com/justgook/goplatformer/pkg/game/system"
	"github.com/justgook/goplatformer/pkg/util"
	"golang.org/x/exp/slog"
)

var _ state.Scene = (*PlayScene)(nil)

type PlayScene struct {
	systems system.Systems


	// OLD STUFF
	currentRoomId int
	Player *system.Player
	LevelManager  *system.LevelManager
}

// Draw implements Scene.
func (world *PlayScene) Draw(screen *ebiten.Image) {
	world.systems.Draw(screen)
}

// Init implements Scene.
func (world *PlayScene) Init() {
	Player := &system.Player{}
	LevelManager := &system.LevelManager{}
	world.systems = system.Systems{
		LevelManager,
		Player,
		&system.UI{},
	}

	world.systems.Init()
	LevelManager.Load(goplatformer.EmbeddedLevel)

	slog.Warn("Move player to LevelManager / collision detection")
	world.LevelManager = LevelManager
	Player.HitExitCallback = world.draftExitsSystem
	world.Player = Player
	LevelManager.Space.Add(Player.Object)
	world.draftExitsSystem(system.ExitEast)
}

// Terminate implements Scene.
func (world *PlayScene) Terminate() {
	world.systems.Terminate()
}

// Update implements Scene.
func (world *PlayScene) Update(aa *state.GameState) error {

	if err := world.systems.Update(aa); err != nil {
		return util.Catch(err)
	}
	return nil
}


func (world *PlayScene) draftExitsSystem(exit system.RoomExit) {
	switch world.currentRoomId {
	case 0:
		switch exit {
		case system.ExitNorth:
			world.currentRoomId = 2
			world.Player.Object.X = 96
			world.Player.Object.Y = 280
		case system.ExitEast:
			world.currentRoomId = 4
			world.Player.Object.X = 32
			world.Player.Object.Y = 32
		case system.ExitSouth:
			world.currentRoomId = 1
			world.Player.Object.X = 228
			world.Player.Object.Y = 60
		case system.ExitWest:
			world.currentRoomId = 3
			world.Player.Object.X = 448
			world.Player.Object.Y = 200
		}
	case 1:
		world.currentRoomId = 0
		world.Player.Object.X = 432
		world.Player.Object.Y = 312
	case 2:
		world.currentRoomId = 0
		world.Player.Object.X = 224
		world.Player.Object.Y = 32
	case 3:
		world.currentRoomId = 0
		world.Player.Object.X = 48
		world.Player.Object.Y = 288
	case 4:
		world.currentRoomId = 0
		world.Player.Object.X = 448
		world.Player.Object.Y = 200
	}
	world.LevelManager.ChangeRoom(world.currentRoomId)
	world.LevelManager.Space.Add(world.Player.Object)
}

