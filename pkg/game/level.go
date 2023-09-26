package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/resolv"
)

type LevelBinary struct {
	Rooms    []*RoomBinary
	tileMaps []*image.Image
}

type RoomBinary struct {
	layers    []int
	collision []*resolv.Object
}

type Room struct {
	layers    []*ebiten.Image
	collision []*resolv.Object
}
