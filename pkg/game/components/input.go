package components

import (
	"github.com/justgook/goplatformer/pkg/game/state"
	"github.com/yohamta/donburi"
)

var Input = donburi.NewComponentType[state.Input]()
