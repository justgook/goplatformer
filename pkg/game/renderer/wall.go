package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/yohamta/donburi/ecs"
)

func Wall(ecs *ecs.ECS, screen *ebiten.Image) {
	spaceEnt, _ := components.Space.First(ecs.World)
	space := components.Space.Get(spaceEnt)
	drawColor := color.RGBA{R: 60, G: 60, B: 60, A: 255}

	for _, obj := range space.Objects() {
		if obj.HaveTags("solid") {
			vector.DrawFilledRect(screen, float32(obj.X), float32(obj.Y), float32(obj.W), float32(obj.H), drawColor, false)
		}
	}
}
