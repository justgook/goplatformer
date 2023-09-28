package stage

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer"
	"github.com/justgook/goplatformer/pkg/bin"
	"github.com/justgook/goplatformer/pkg/game/system"
	"github.com/justgook/goplatformer/pkg/util"
	"image"
	"image/color"
	_ "image/png"
)

type Play struct {
	Level   *bin.Level
	TileSet *ebiten.Image

	Player        *system.Player
	Room          *system.Room
	currentRoomId int
}

func (world *Play) Init() {
	world.Level = &bin.Level{}

	util.OrDie(world.Level.Load(goplatformer.EmbeddedLevel))
	img, _ := util.Get2OrDie(image.Decode(bytes.NewReader(world.Level.Image)))
	world.TileSet = ebiten.NewImageFromImage(img)

	world.Player = system.NewPlayer(world.draftExits)
	util.OrDie(world.Player.Animation.Load(goplatformer.EmbeddedPlayerSprite))

	world.draftExits(system.ExitEast)
}

func (world *Play) Update() {
	world.Room.Update()
	world.Player.Update()

}

func (world *Play) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{G: 10, B: 10, A: 255})
	/* ===================================================== */
	screen.DrawImage(world.Room.Draw(), nil)
	/* ===================================================== */

	/* ===================================================== */
	/* Player Sprite*/
	player := world.Player.Object
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(player.X)-16, float64(player.Y)-16)
	screen.DrawImage(world.Player.Draw(), op)
	/* ===================================================== */
}

func (world *Play) draftExits(exit system.RoomExit) {
	//
	switch world.currentRoomId {
	case 0:
		switch exit {
		case system.ExitNorth:
			world.currentRoomId = 2
			world.Player.Object.X = 100
			world.Player.Object.Y = 420
		case system.ExitEast:
			world.currentRoomId = 4
			world.Player.Object.X = 32
			world.Player.Object.Y = 32
		case system.ExitSouth:
			world.currentRoomId = 1
			world.Player.Object.X = 228
			world.Player.Object.Y = 60
		case system.ExitWest:
			world.currentRoomId = 3
			world.Player.Object.X = 448
			world.Player.Object.Y = 200
		}
	case 1:
		world.currentRoomId = 0
		world.Player.Object.X = 432
		world.Player.Object.Y = 448
	case 2:
		world.currentRoomId = 0
		world.Player.Object.X = 224
		world.Player.Object.Y = 32
	case 3:
		world.currentRoomId = 0
		world.Player.Object.X = 48
		world.Player.Object.Y = 288
	case 4:
		world.currentRoomId = 0
		world.Player.Object.X = 448
		world.Player.Object.Y = 200
	}
	world.Room = system.NewRoom(world.TileSet, 16, world.Level.Rooms[world.currentRoomId])
	world.Room.Space.Add(world.Player.Object)
}
