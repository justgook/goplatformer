package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/tags"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func Ramp(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Ramp.Each(ecs.World, func(e *donburi.Entry) {
		o := components.Object.Get(e)
		drawColor := color.RGBA{R: 255, G: 50, B: 100, A: 255}
		tri := o.Shape.(*resolv.ConvexPolygon)
		drawPolygon(screen, tri, drawColor)
	})
}

func drawPolygon(screen *ebiten.Image, polygon *resolv.ConvexPolygon, color color.Color) {
	for _, line := range polygon.Lines() {
		vector.StrokeLine(
			screen,
			float32(line.Start.X()),
			float32(line.Start.Y()),
			float32(line.End.X()),
			float32(line.End.Y()),
			1,
			color,
			false,
		)

	}
}
