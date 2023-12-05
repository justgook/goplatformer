package system

import (
	"log/slog"

	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateBulletEmitter(ecs *ecs.ECS) {
	components.BulletEmitter.Each(ecs.World, func(e *donburi.Entry) {
		bullet := components.BulletEmitter.Get(e)
		err := bullet.Runner.Update()
		if err != nil {
			slog.Error("UpdateBulletEmitter", "error", util.Catch(err))
		}

	})
}

func UpdateBullets(ecs *ecs.ECS) {
	components.Bullet.Each(ecs.World, func(e *donburi.Entry) {
		b := components.Bullet.Get(e).Runner

		if err := b.Update(); err != nil || b.Vanished() {
			ecs.World.Remove(e.Entity())
		}

	})
}
