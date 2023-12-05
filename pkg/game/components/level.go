package components

import (
	"github.com/justgook/goplatformer/pkg/game/components/level"
	"github.com/yohamta/donburi"
)

var Level = donburi.NewComponentType[level.Data]()
