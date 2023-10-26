package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/yohamta/donburi/ecs"
)

func Level(ecs *ecs.ECS, screen *ebiten.Image) {
	entity, ok := components.Level.First(ecs.World)
	if !ok {
		return
	}

	level := components.Level.Get(entity)
	for _, room := range level.RoomsResults {
		for _, layer := range room.Render {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(room.Rect.Min.X), float64(room.Rect.Min.Y))
			screen.DrawImage(layer, op)
		}
	}
}
