package resources

import (
	"github.com/justgook/goplatformer/pkg/resolv/v2"
)

type TagType = int64
type Level struct {
	Rooms []*Room
	Image []byte
}

type Doors struct {
	N bool
	E bool
	S bool
	W bool
}

type Tile struct {
	X int64
	Y int64
	T int64
}

type Room struct {
	Layers    [][]Tile
	Doors     Doors
	W         int
	H         int
	Collision []*resolv.Object[TagType]
}

