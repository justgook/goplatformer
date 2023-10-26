package entity

import (
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSpace(ecs *ecs.ECS) *donburi.Entry {
	space := Space.Spawn(ecs)

	spaceData := resolv.NewSpace[string](1, 1, 16, 16)
	components.Space.Set(space, spaceData)

	return space
}
