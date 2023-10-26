package components

import (
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/yohamta/donburi"
)

var Object = donburi.NewComponentType[resolv.Object[string]]()
