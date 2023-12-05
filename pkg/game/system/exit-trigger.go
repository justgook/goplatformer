package system

import (
	. "github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/yohamta/donburi/ecs"
	"image"
	"log/slog"
)

func ExitTrigger(ecs *ecs.ECS) {
	entity, ok := components.Level.First(ecs.World)
	if !ok {
		return
	}

	levelComp := components.Level.Get(entity)

	spaceEnt, _ := components.Space.First(ecs.World)
	space := components.Space.Get(spaceEnt)

	playerEntry, _ := components.Player.First(ecs.World)
	player := components.Player.Get(playerEntry)
	playerObject := components.Object.Get(playerEntry)

	if check := playerObject.Check(0, 0,
		ObjectTagExitTriggerNorth,
		ObjectTagExitTriggerEast,
		ObjectTagExitTriggerSouth,
		ObjectTagExitTriggerWest,
	); check != nil {
		nextP := levelComp.CurrentRoomXY
		switch {
		case check.HasTags(ObjectTagExitTriggerNorth):
			nextP = nextP.Add(image.Pt(0, -1))
		case check.HasTags(ObjectTagExitTriggerEast):
			nextP = nextP.Add(image.Pt(1, 0))
		case check.HasTags(ObjectTagExitTriggerSouth):
			nextP = nextP.Add(image.Pt(0, 1))
		case check.HasTags(ObjectTagExitTriggerWest):
			nextP = nextP.Add(image.Pt(-1, 0))
		}
		slog.Info("system.ExitTrigger", "check", nextP)
		levelComp.SetCurrent(nextP, space)
		exitP := image.Point{}
		switch {
		case check.HasTags(ObjectTagExitTriggerNorth):
			exitP = levelComp.CurrentRoom.LevelEnter.EnterS
		case check.HasTags(ObjectTagExitTriggerEast):
			exitP = levelComp.CurrentRoom.LevelEnter.EnterW
		case check.HasTags(ObjectTagExitTriggerSouth):
			exitP = levelComp.CurrentRoom.LevelEnter.EnterN
		case check.HasTags(ObjectTagExitTriggerWest):
			exitP = levelComp.CurrentRoom.LevelEnter.EnterE
		}

		playerObject.X = float64(exitP.X) - playerObject.W*0.5
		playerObject.Y = float64(exitP.Y) - playerObject.H
		player.SpeedX = 0
		player.SpeedY = 0
		player.OnGround = nil
		player.WallSliding = nil
		player.IgnorePlatform = nil
		space.Add(playerObject)
	}
}
