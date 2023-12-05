package components

import (
	"github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/yohamta/donburi"
)

var Space = donburi.NewComponentType[domain.ObjectSpace]()
