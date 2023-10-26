package system

import (
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateFloatingPlatform(ecs *ecs.ECS) {
	tags.FloatingPlatform.Each(ecs.World, func(e *donburi.Entry) {
		tw := components.Tween.Get(e)
		// Platform movement needs to be done first to make sure there's no space between the top and the player's bottom; otherwise, an alternative might
		// be to have the platform detect to see if the Player's resting on it, and if so, move the player up manually.
		y, _, seqDone := tw.Update(1.0 / 60.0)

		obj := components.Object.Get(e)
		obj.Y = float64(y)
		if seqDone {
			tw.Reset()
		}
	})
}
