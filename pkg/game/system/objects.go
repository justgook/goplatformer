package system

import (
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateObjects(ecs *ecs.ECS) {
	components.Object.Each(ecs.World, func(e *donburi.Entry) {
		components.Object.Get(e).Update()
	})
}
