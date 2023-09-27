package stage

import (
	"bytes"
	"image"
	_ "image/png"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/justgook/goplatformer"
	"github.com/justgook/goplatformer/pkg/bin"
	"github.com/justgook/goplatformer/pkg/game/system"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/justgook/goplatformer/pkg/util"
)

type Play struct {
	Space   *resolv.Space[bin.TagType]
	Level   *bin.Level
	TileSet *ebiten.Image
	Player  *system.Player
}

func (world *Play) Init() {
	world.Level = &bin.Level{}

	//// INITING LEVEL _ TODO MEOVE ME OUTSIDE
	util.OrDie(world.Level.Load(goplatformer.EmbeddedLevel))
	img, _ := util.Get2OrDie(image.Decode(bytes.NewReader(world.Level.Image)))
	world.TileSet = ebiten.NewImageFromImage(img)

	room := world.Level.Rooms[0]
	// Define the world's Space.
	// Here, a Space is essentially a grid (the game's width and height),
	// made up of 16x16 cells. Each cell can have 0 or more Objects within it,
	// and collisions can be found by checking the Space
	// to see if the Cells at specific positions contain (or would contain) Objects.
	// This is a broad, simplified approach to collision detection.
	world.Space = resolv.NewSpace[bin.TagType](room.W, room.H, 16, 16)

	// Construct the solid level geometry.
	// Note that the simple approach of checking cells in a Space for collision works simply when the geometry is
	// aligned with the cells, as it all is in this platformer example.
	world.Space.Add(room.Collision...)

	world.Player = system.NewPlayer(world.Space)

	util.OrDie(world.Player.Animation.Load(goplatformer.EmbeddedPlayerSprite))
}

func (world *Play) Update() {
	system.PlayerUpdate(world.Player)
}

func (world *Play) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{G: 10, B: 10, A: 255})
	/* ===================================================== */
	world.DrawLevel(screen)
	/* ===================================================== */

	/* ===================================================== */
	/* Player Sprite*/
	player := world.Player.Object
	op := &ebiten.DrawImageOptions{}
	if !world.Player.FacingRight {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(48, 0)
	}
	op.GeoM.Translate(float64(player.X)-16, float64(player.Y)-16)
	screen.DrawImage(world.Player.Animation.Sprite, op)
	/* ===================================================== */
}

func (world *Play) DrawDebug(screen *ebiten.Image) {
	drawColor := color.RGBA{R: 255, G: 50, B: 100, A: 255}
	for _, o := range world.Space.Objects() {
		if o.HaveTags(3) {
			drawColor = color.RGBA{R: 255, G: 50, B: 100, A: 255}
		} else if o.HaveTags(1) {
			drawColor = color.RGBA{R: 50, G: 255, B: 100, A: 255}
		} else if o.HaveTags(99) /*Player*/ {
			continue
		} else {
			drawColor = color.RGBA{R: 50, G: 50, B: 255, A: 255}
		}
		vector.StrokeRect(screen, float32(o.X+1), float32(o.Y+1), float32(o.W-2), float32(o.H-2), 1, drawColor, false)
	}
}

func (world *Play) DrawLevel(screen *ebiten.Image) {
	room := world.Level.Rooms[0]
	for _, layer := range room.Layers {
		for _, tile := range layer {
			id := tile.T
			x1 := int(id%12) * 16
			y1 := int(id/12) * 16
			rect := image.Rect(x1, y1, x1+16, y1+16)
			// TODO prebake all tiles instead of using it here, and just use direct index
			result := world.TileSet.SubImage(rect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.X), float64(tile.Y))
			screen.DrawImage(result, op)
		}
	}
}
