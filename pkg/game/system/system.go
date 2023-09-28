package system

import "github.com/hajimehoshi/ebiten/v2"

type System interface {
	Init()
	Terminate()
	Update()
	Draw() *ebiten.Image
}

type RoomExit = int64

const (
	ExitNorth RoomExit = 5
	ExitEast  RoomExit = 6
	ExitSouth RoomExit = 7
	ExitWest  RoomExit = 8
)
