package components

import (
	"github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/yohamta/donburi"
)

type PlayerData struct {
	SpeedX         float64
	SpeedY         float64
	OnGround       *domain.Object
	WallSliding    *domain.Object
	FacingRight    bool
	IgnorePlatform *domain.Object
}

var Player = donburi.NewComponentType[PlayerData]()
