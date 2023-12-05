package renderer

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var bulletImg = func() *ebiten.Image {
	img := ebiten.NewImage(6, 6)
	vector.DrawFilledCircle(img, 3, 3, 3, color.RGBA{R: 0xe8, G: 0x7a, B: 0x90, A: 0xff}, true)
	return img
}()

func Bullets(ecs *ecs.ECS, screen *ebiten.Image) {
	count := 0
	components.Bullet.Each(ecs.World, func(e *donburi.Entry) {
		b := components.Bullet.Get(e).Runner
		//for i := range bullet.Bullets {
		count++
		x, y := b.Position()
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(x-3, y-3)
		screen.DrawImage(bulletImg, opts)
		//}
	})

	//DEBUG INFO
	playerEntry, _ := components.Player.First(ecs.World)
	playerObject := components.Object.Get(playerEntry)
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Bullets: %d", count),
		int(playerObject.X)-36,
		int(playerObject.Y)-16,
	)
}
