package components

import (
	"github.com/justgook/goplatformer/pkg/bulletml"
	"github.com/yohamta/donburi"
)

type BulletEmitterData struct {
	Runner bulletml.Runner
}

type BulletData struct {
	Runner bulletml.BulletRunner
}

var BulletEmitter = donburi.NewComponentType[BulletEmitterData]()
var Bullet = donburi.NewComponentType[BulletData]()
