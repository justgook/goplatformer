//go:build production

package goplatformer

import (
	_ "embed"
)

//go:embed build/sprite/Player.sprite
var EmbeddedPlayerSprite []byte

//go:embed build/level/level1.level
var EmbeddedLevel []byte
