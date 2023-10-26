package components

import (
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/yohamta/donburi"
)

var Space = donburi.NewComponentType[resolv.Space[string]]()
