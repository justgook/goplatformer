package game

import "github.com/hajimehoshi/ebiten/v2"

type Stage interface {
	Init()
	Update()
	Draw(*ebiten.Image)
}
