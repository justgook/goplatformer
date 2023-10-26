package system

import (
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateCharAnim(ecs *ecs.ECS) {
	components.CharAnim.Each(ecs.World, func(e *donburi.Entry) {
		components.CharAnim.Get(e).Update()
	})
}
