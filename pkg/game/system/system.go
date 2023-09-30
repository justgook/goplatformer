package system

import "github.com/hajimehoshi/ebiten/v2"

type System interface {
	Init()
	Terminate()
	Update()
	Draw(screen *ebiten.Image)
}
type Systems []System

func (s Systems) Init() {
	for i := range s {
		s[i].Init()
	}
}

func (s Systems) Draw(screen *ebiten.Image) {
	for i := range s {
		s[i].Draw(screen)
	}
}

func (s Systems) Update() {
	for i := range s {
		s[i].Update()
	}
}

func (s Systems) Terminate() {
	for i := range s {
		s[i].Terminate()
	}
}

type RoomExit = int64

const (
	ExitNorth RoomExit = 5
	ExitEast  RoomExit = 6
	ExitSouth RoomExit = 7
	ExitWest  RoomExit = 8
)
