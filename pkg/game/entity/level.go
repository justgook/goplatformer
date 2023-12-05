package entity

import (
	"github.com/justgook/goplatformer"
	. "github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/components/level"
	"github.com/justgook/goplatformer/pkg/game/state"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var Level = newArchetype(
	components.Level,
)

func CreateLevel(ecs *ecs.ECS, st *state.GameState, collisionSpace *ObjectSpace) *donburi.Entry {
	output := Level.Spawn(ecs)
	goalDistance, branchLength := 6, 3

	newLevel := level.New(
		st.Rand,
		goalDistance,
		branchLength,
		goplatformer.EmbeddedLevel,
		collisionSpace,
	)

	components.Level.Set(output, newLevel)

	return output
}
