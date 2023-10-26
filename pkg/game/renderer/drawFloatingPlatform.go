package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"image/color"
)

func FloatingPlatform(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.FloatingPlatform.Each(ecs.World, func(e *donburi.Entry) {
		o := components.Object.Get(e)
		drawColor := color.RGBA{R: 180, G: 100, A: 255}
		vector.DrawFilledRect(screen, float32(o.X), float32(o.Y), float32(o.W), float32(o.H), drawColor, false)
	})
}
