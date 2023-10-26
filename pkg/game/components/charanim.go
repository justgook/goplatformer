package components

import (
	"github.com/justgook/goplatformer/pkg/game/components/sprite"
	"github.com/yohamta/donburi"
)

type CharacterData struct {
	Run          *sprite.Animation
	JumpUp       *sprite.Animation
	JumpMax      *sprite.Animation
	Fall         *sprite.Animation
	WallSlideLow *sprite.Animation
	Idle         *sprite.Animation

	Current *sprite.Animation
}

var CharAnim = donburi.NewComponentType[CharacterData]()

func (c *CharacterData) Update() {
	c.Current.Update()
}
