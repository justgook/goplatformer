package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/components/sprite"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func CharAnim(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(filter.Contains(components.CharAnim, components.Object))
	query.Each(ecs.World, func(e *donburi.Entry) {
		o := components.Object.Get(e)
		s := components.CharAnim.Get(e).Current
		// ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
		getFromSpriteSizeFW := float64(s.W())
		getFromSpriteSizeFH := float64(s.H())
		x := o.X + (o.W-getFromSpriteSizeFW)*0.5
		y := o.Y + (o.H-getFromSpriteSizeFH)*0.5
		s.Draw(screen, sprite.DrawOpts(x, y))
	})

}
