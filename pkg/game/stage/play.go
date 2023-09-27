package stage

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/justgook/goplatformer"
	"github.com/justgook/goplatformer/pkg/bin"
	"github.com/justgook/goplatformer/pkg/game/system"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/justgook/goplatformer/pkg/util"
	"image/color"
)

type Play struct {
	Space  *resolv.Space[bin.TagType]
	Level  *bin.Level
	Player *system.Player
}

func (world *Play) Init() {
	world.Level = &bin.Level{}
	util.OrDie(world.Level.Load(goplatformer.EmbeddedLevel))

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
	screen.Fill(color.RGBA{G: 145, B: 255, A: 255})
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
