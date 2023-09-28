package stage

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/justgook/goplatformer"
	"github.com/justgook/goplatformer/pkg/bin"
	"github.com/justgook/goplatformer/pkg/game/system"
	"github.com/justgook/goplatformer/pkg/game/ui"
	"github.com/justgook/goplatformer/pkg/util"
)

type Play struct {
	UI     *ebitenui.UI
	HUD    *widget.Container
	MAP_UI *widget.Container

	Level   *bin.Level
	TileSet *ebiten.Image

	Player        *system.Player
	Room          *system.Room
	currentRoomId int
}

func (world *Play) Init() {
	bitmapFont := util.GetOrDie(util.LoadFont(goplatformer.ExcelFont, 16))
	world.HUD = ui.NewHUD(bitmapFont)
	world.MAP_UI = ui.NewMap(bitmapFont)
	world.HUD.AddChild(world.MAP_UI)
	world.UI = &ebitenui.UI{
		Container: world.HUD,
	}
	// ====================================================================================

	world.Level = &bin.Level{}

	util.OrDie(world.Level.Load(goplatformer.EmbeddedLevel))
	img, _ := util.Get2OrDie(image.Decode(bytes.NewReader(world.Level.Image)))
	world.TileSet = ebiten.NewImageFromImage(img)

	world.Player = system.NewPlayer(world.draftExitsSystem)
	util.OrDie(world.Player.Animation.Load(goplatformer.EmbeddedPlayerSprite))

	world.draftExitsSystem(system.ExitEast)
}

func (world *Play) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyTab) {
    	world.MAP_UI.GetWidget().Visibility = widget.Visibility_Show

		} else {
	world.MAP_UI.GetWidget().Visibility = widget.Visibility_Hide

  }

	world.Room.Update()
	world.Player.Update()
	world.UI.Update()
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
	world.UI.Draw(screen)

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))
}

func (world *Play) draftExitsSystem(exit system.RoomExit) {
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
