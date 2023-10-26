package goplatformer

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/util"
)

var (
	//go:embed asset/intro/logo-min.png
	introLogo0x069 []byte
	IntroLogo0x069 = func() *ebiten.Image {
		return util.GetOrDie(util.ImageFromBytes(introLogo0x069))
	}()
	IntroBg     = "051726"
	IntroAccent = "81ffd9"
)
