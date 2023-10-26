package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func Player(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Player.Each(ecs.World, func(e *donburi.Entry) {
		player := components.Player.Get(e)
		o := components.Object.Get(e)
		drawColor := color.RGBA{G: 255, B: 60, A: 255}
		if player.OnGround == nil {
			// We draw the player as a different color when jumping so we can visually see when he's in the air.
			drawColor = color.RGBA{R: 200, B: 200, A: 255}
		}
		vector.DrawFilledRect(screen, float32(o.X), float32(o.Y), float32(o.W), float32(o.H), drawColor, false)
	})
}
