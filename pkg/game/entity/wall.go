package entity

import (
	"github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateWall(ecs *ecs.ECS, obj *domain.Object) *donburi.Entry {
	wall := Wall.Spawn(ecs)
	components.Object.Set(wall, obj)
	return wall
}
