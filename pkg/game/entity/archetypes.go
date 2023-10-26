package entity

import (
	components2 "github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/layers"
	"github.com/justgook/goplatformer/pkg/game/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var (
	Platform = newArchetype(
		tags.Platform,
		components2.Object,
	)
	FloatingPlatform = newArchetype(
		tags.FloatingPlatform,
		components2.Object,
		components2.Tween,
	)
	Player = newArchetype(
		tags.Player,
		components2.Input,
		components2.CharAnim,
		components2.Player,
		components2.Object,
	)
	Ramp = newArchetype(
		tags.Ramp,
		components2.Object,
	)
	Space = newArchetype(
		components2.Space,
	)
	Wall = newArchetype(
		tags.Wall,
		components2.Object,
	)
)

type archetype struct {
	components []donburi.IComponentType
}

func newArchetype(cs ...donburi.IComponentType) *archetype {
	return &archetype{
		components: cs,
	}
}

func (a *archetype) Spawn(ecs *ecs.ECS, cs ...donburi.IComponentType) *donburi.Entry {
	e := ecs.World.Entry(ecs.Create(
		layers.Background,
		append(a.components, cs...)...,
	))
	return e
}
