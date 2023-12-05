package system

import (
	"github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/yohamta/donburi/ecs"
	"log/slog"
)

func EnemyTrigger(ecs *ecs.ECS) {
	entity, ok := components.Level.First(ecs.World)
	if !ok {
		return
	}

	level := components.Level.Get(entity)
	_ = level

	playerEntry, _ := components.Player.First(ecs.World)
	player := components.Player.Get(playerEntry)
	playerObject := components.Object.Get(playerEntry)
	_, _ = player, playerObject

	if check := playerObject.Check(0, 0, domain.ObjectTagEnemyTrigger); check != nil {
		slog.Info("system.EnemyTrigger", "check", check)
	}
	//else {
	//	slog.Info("system.EnemyTrigger", "check", check)
	//}
}
