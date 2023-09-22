package main

import "github.com/hajimehoshi/ebiten/v2"

type StageInterface interface {
	Init()
	Update()
	Draw(*ebiten.Image)
}
