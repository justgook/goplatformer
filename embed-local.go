//go:build !production

package goplatformer

import (
	_ "embed"
)

//go:embed build.nosync/sprite/Player.char
var playerAnim []byte

//go:embed build.nosync/level/level1.level
var embeddedLevel []byte
