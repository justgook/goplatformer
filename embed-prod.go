//go:build production

package goplatformer

import (
	_ "embed"
)

//go:embed build/sprite/Player.char
var playerAnim []byte

//go:embed build/level/level1.level
var embeddedLevel []byte
