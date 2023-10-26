package components

import (
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/yohamta/donburi"
)

type PlayerData struct {
	SpeedX         float64
	SpeedY         float64
	OnGround       *resolv.Object[string]
	WallSliding    *resolv.Object[string]
	FacingRight    bool
	IgnorePlatform *resolv.Object[string]
}

var Player = donburi.NewComponentType[PlayerData]()
