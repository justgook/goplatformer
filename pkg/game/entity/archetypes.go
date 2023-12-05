package entity

import (
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/layers"
	"github.com/justgook/goplatformer/pkg/game/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var (
	Platform = newArchetype(
		tags.Platform,
		components.Object,
	)
	FloatingPlatform = newArchetype(
		tags.FloatingPlatform,
		components.Object,
		components.Tween,
	)
	Player = newArchetype(
		tags.Player,
		components.Input,
		components.CharAnim,
		components.Player,
		components.Object,
	)
	Ramp = newArchetype(
		tags.Ramp,
		components.Object,
	)
	Space = newArchetype(
		components.Space,
	)
	Wall = newArchetype(
		tags.Wall,
		components.Object,
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
