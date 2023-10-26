package entity

import (
	components2 "github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreatePlatform(ecs *ecs.ECS, object *resolv.Object[string]) *donburi.Entry {
	platform := Platform.Spawn(ecs)
	components2.Object.Set(platform, object)

	return platform
}

func CreateFloatingPlatform(ecs *ecs.ECS, object *resolv.Object[string]) *donburi.Entry {
	platform := FloatingPlatform.Spawn(ecs)
	components2.Object.Set(platform, object)

	// The floating platform moves using a *gween.Sequence sequence of tweens, moving it back and forth.
	tw := gween.NewSequence()
	obj := components2.Object.Get(platform)
	tw.Add(
		gween.New(float32(obj.Y), float32(obj.Y-128), 2, ease.Linear),
		gween.New(float32(obj.Y-128), float32(obj.Y), 2, ease.Linear),
	)
	components2.Tween.Set(platform, tw)

	return platform
}
