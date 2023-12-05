package entity

import (
	"github.com/justgook/goplatformer/pkg/bulletml"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var BulletEmitter = newArchetype(
	components.BulletEmitter,
	// components.Object,
)
var Bullet = newArchetype(
	components.Bullet,
	// components.Object,
)

func CreateBulletEmitter(
	ecs *ecs.ECS,
	bml *bulletml.BulletML,
	target donburi.Entity,
	emitter donburi.Entity,
) *donburi.Entry {
	bulletEmitter := BulletEmitter.Spawn(ecs)

	opts := &bulletml.NewRunnerOptions{
		// Called when new bullet fired
		OnBulletFired: func(bulletRunner bulletml.BulletRunner, _ *bulletml.FireContext) {
			bullet := Bullet.Spawn(ecs)
			components.Bullet.Set(bullet, &components.BulletData{Runner: bulletRunner})
		},
		// Tell current emitter position
		CurrentShootPosition: func() (float64, float64) {
			entity := ecs.World.Entry(emitter)
			obj := components.Object.Get(entity)

			return obj.X, obj.Y
		},

		// Tell current target position
		CurrentTargetPosition: func() (float64, float64) {
			entity := ecs.World.Entry(target)
			obj := components.Object.Get(entity)
			return obj.X, obj.Y
		},
	}
	runner, _ := bulletml.NewRunner(bml, opts)
	bulletEmitterData := &components.BulletEmitterData{Runner: runner}

	components.BulletEmitter.Set(bulletEmitter, bulletEmitterData)
	// components.Object.Set(bulletEmitter, object)
	return bulletEmitter
}
