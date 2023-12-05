package entity

import (
	"github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/state"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreatePlayer(ecs *ecs.ECS, anim components.CharacterData, input *state.Input) *donburi.Entry {
	player := Player.Spawn(ecs)

	obj := resolv.NewObject[domain.ObjectTag](32, 128, 16, 24)
	components.Object.Set(player, obj)
	components.Player.SetValue(player, components.PlayerData{
		FacingRight: true,
	})
	obj.SetShape(resolv.NewRectangle(0, 0, 16, 24))
	components.CharAnim.SetValue(player, anim)
	components.Input.SetValue(player, *input)

	return player
}
