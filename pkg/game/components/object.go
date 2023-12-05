package components

import (
	"github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/yohamta/donburi"
)

var Object = donburi.NewComponentType[domain.Object]()
