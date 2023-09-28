package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/bin"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"image"
	"image/color"
)

type Room struct {
	Space   *resolv.Space[bin.TagType]
	TileSet *ebiten.Image
	render  *ebiten.Image
}

func (r *Room) Init() {

}

func (r *Room) Terminate() {

}

func (r *Room) Update() {

}

func (r *Room) Draw() *ebiten.Image {
	return r.render
}

func (r *Room) bakeImage(source *bin.Room) {
	r.render.Fill(color.RGBA{})
	for _, layer := range source.Layers {
		for _, tile := range layer {
			id := tile.T
			x1 := int(id%12) * 16
			y1 := int(id/12) * 16
			rect := image.Rect(x1, y1, x1+16, y1+16)
			result := r.TileSet.SubImage(rect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.X), float64(tile.Y))
			r.render.DrawImage(result, op)
		}
	}
}

func NewRoom(tileSet *ebiten.Image, cellSize int, input *bin.Room) *Room {
	output := &Room{
		TileSet: tileSet,
		Space:   resolv.NewSpace[bin.TagType](input.W, input.H, cellSize, cellSize),
		render:  ebiten.NewImage(input.W, input.H),
	}

	output.Space.Add(input.Collision...)
	output.bakeImage(input)

	return output
}
