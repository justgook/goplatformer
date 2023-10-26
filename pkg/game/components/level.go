package components

import (
	"github.com/justgook/goplatformer/pkg/game/components/level"
	"github.com/yohamta/donburi"
)

const (
	ExitN = level.ExitN
	ExitE = level.ExitE
	ExitS = level.ExitS
	ExitW = level.ExitW
)

var Level = donburi.NewComponentType[level.Data]()
