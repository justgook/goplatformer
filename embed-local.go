//go:build !production

package goplatformer

import (
	_ "embed"
)

//go:embed build.nosync/sprite/Player.sprite
var EmbeddedPlayerSprite []byte

//go:embed build.nosync/level/level1.level
var EmbeddedLevel []byte

	// util.OrDie(world.currentLevel.Load(goplatformer.EmbeddedLevel))

