package entity

import (
	components2 "github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/state"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreatePlayer(ecs *ecs.ECS, anim components2.CharacterData, input *state.Input) *donburi.Entry {
	player := Player.Spawn(ecs)

	obj := resolv.NewObject[string](32, 128, 16, 24)
	components2.Object.Set(player, obj)
	components2.Player.SetValue(player, components2.PlayerData{
		FacingRight: true,
	})
	obj.SetShape(resolv.NewRectangle(0, 0, 16, 24))
	components2.CharAnim.SetValue(player, anim)
	components2.Input.SetValue(player, *input)

	return player
}
