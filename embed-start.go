package goplatformer

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/util"
	_ "image/png"
)

var (
	//go:embed asset/start/title.png
	startMenuTitleImg []byte
	StartMenuTitleImg  = func () *ebiten.Image {
		return util.GetOrDie(util.LoadImage(startMenuTitleImg))
	}()
)

