package entity

import (
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateWall(ecs *ecs.ECS, obj *resolv.Object[string]) *donburi.Entry {
	wall := Wall.Spawn(ecs)
	components.Object.Set(wall, obj)
	return wall
}
