package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/components/level"
	"github.com/justgook/goplatformer/pkg/ui"
	"github.com/yohamta/donburi/ecs"
	"image"
	"image/color"
)

func Level(ecs *ecs.ECS, screen *ebiten.Image) {
	entity, ok := components.Level.First(ecs.World)
	if !ok {
		return
	}

	levelComp := components.Level.Get(entity)

	for _, layer := range levelComp.CurrentRoom.RenderLayer {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(layer, op)
	}

	/// DEBUG

	//spaceEnt, _ := components.Space.First(ecs.World)
	//space := components.Space.Get(spaceEnt)
	//triggerEnemy(levelComp, screen)
	//triggerExits(space, screen)
	//enterPoints(levelComp, screen)
}

func enterPoints(info *level.Data, screen *ebiten.Image) {
	for _, a := range []image.Point{
		info.CurrentRoom.LevelEnter.EnterN,
		info.CurrentRoom.LevelEnter.EnterE,
		info.CurrentRoom.LevelEnter.EnterS,
		info.CurrentRoom.LevelEnter.EnterW,
	} {
		ui.FillRect(screen, &ui.FillRectOpts{
			Rect:  image.Rect(a.X-1, a.Y-1, a.X+1, a.Y+1),
			Color: color.RGBA{R: 0xff, A: 0xff},
		})
	}

}

func triggerExits(space *domain.ObjectSpace, screen *ebiten.Image) {
	for _, obj := range space.Objects() {
		if obj.HaveTags(
			domain.ObjectTagExitTriggerNorth,
			domain.ObjectTagExitTriggerEast,
			domain.ObjectTagExitTriggerSouth,
			domain.ObjectTagExitTriggerWest,
		) {
			renderExit(obj.X, obj.Y, obj.W, obj.H, screen)
		}
	}
}

func renderExit(x, y, w, h float64, screen *ebiten.Image) {
	fillColor := color.RGBA{R: 0xFF, G: 0x00, B: 0xcc, A: 0x12}
	lineColor := color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xff}
	v := image.Rect(
		int(x),
		int(y),
		int(x)+int(w),
		int(y)+int(h),
	)
	ui.FillRect(
		screen,
		&ui.FillRectOpts{
			Rect:  v,
			Color: fillColor,
		},
	)
	ui.DrawRect(
		screen,
		&ui.DrawRectOpts{
			Rect:        v,
			Color:       lineColor,
			StrokeWidth: 2,
		},
	)
}
func triggerEnemy(level *level.Data, screen *ebiten.Image) {
	fillColor := color.RGBA{R: 0xFF, G: 0xcc, B: 0x00, A: 0x12}
	lineColor := color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xff}
	for _, v := range level.CurrentRoom.TriggerSpawnEnemy {
		ui.FillRect(
			screen,
			&ui.FillRectOpts{
				Rect:  v.Area,
				Color: fillColor,
			},
		)
		ui.DrawRect(
			screen,
			&ui.DrawRectOpts{
				Rect:        v.Area,
				Color:       lineColor,
				StrokeWidth: 2,
			},
		)
	}
}
