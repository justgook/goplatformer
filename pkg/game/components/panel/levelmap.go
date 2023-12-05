package panel

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/game/components/level"
	"github.com/justgook/goplatformer/pkg/ui"
)

type LevelMap struct {
	View *ui.View
}

func (l *LevelMap) Init(data *level.Data, currentRoom *image.Point) {
	l.View = &ui.View{
		WidthInPct:  100,
		HeightInPct: 100,
		Justify:     ui.JustifyStart,
		AlignItems:  ui.AlignItemEnd,
		// MarginBottom: 32, // TODO find out why it not works
		// That is just hack because MarginBottom not works
		Position: ui.PositionAbsolute,
		Bottom:   ui.Int(32),
	}

	wrapper := &ui.View{
		//Handler: &Box{
		//	Color: color.RGBA{R: 0xff, G: 0xcc, B: 0xcc, A: 0xff},
		//},
	}
	margin := 4
	roomW := 16
	roomH := 16
	finalW := 1
	finalH := 1
	for k, item := range data.Rooms {
		top := margin + (roomH+margin)*k.Y
		left := margin + (roomW+margin)*k.X
		wrapper.AddChild(&ui.View{
			Position: ui.PositionAbsolute,
			Width:    roomW,
			Height:   roomH,
			Top:      top,
			Left:     left,
			Handler: &MapRoomUI{
				Features:    item.Features,
				Point:       k,
				CurrentRoom: currentRoom,
			},
		})
		finalW = max(finalW, left+margin+roomW)
		finalH = max(finalH, top+margin+roomH)
	}
	wrapper.Width = finalW
	wrapper.Height = finalH
	l.View.AddChild(wrapper)
}

type MapRoomUI struct {
	Features RoomNavigation
	image.Point
	//
	CurrentRoom *image.Point
}

func (mr *MapRoomUI) Draw(screen *ebiten.Image, frame image.Rectangle, v *ui.View) {
	cellColor := color.RGBA{R: 0x77, G: 0x66, B: 0x66, A: 0xff}
	if mr.Features.Has(RoomNavigationStart) {
		cellColor.G = 0xff
	} else if mr.Features.Has(RoomNavigationGoal) {
		cellColor.R = 0xff
	}

	ui.FillRect(screen, &ui.FillRectOpts{
		Rect:  frame,
		Color: cellColor,
	})

	if *mr.CurrentRoom == mr.Point {
		ui.DrawRect(screen, &ui.DrawRectOpts{
			Rect:        frame,
			Color:       color.RGBA{R: 0xFF, A: 0xff},
			StrokeWidth: 2,
		})
	}
	x := frame.Min.X
	y := frame.Min.Y
	w := frame.Dx()
	h := frame.Dy()
	dw := 1
	dh := 2
	doorColor := color.RGBA{R: 0xcc, G: 0xcc, B: 0xcc, A: 0xff}
	if mr.Features.Has(RoomNavigationExitN) {
		ui.FillRect(screen,
			&ui.FillRectOpts{
				Rect:  image.Rect(x-dw+w/2, y-dh, x+dw+w/2, y+dh),
				Color: doorColor,
			},
		)
	}
	if mr.Features.Has(RoomNavigationExitE) {
		ui.FillRect(screen,
			&ui.FillRectOpts{

				Rect:  image.Rect(x-dh+w, y+h/2-dw, x+dh+w, y+h/2+dw),
				Color: doorColor,
			},
		)
	}
	if mr.Features.Has(RoomNavigationExitS) {
		ui.FillRect(screen,
			&ui.FillRectOpts{
				Rect:  image.Rect(x-dw+w/2, y-dh+h, x+dw+w/2, y+dh+h),
				Color: doorColor,
			},
		)
	}
	if mr.Features.Has(RoomNavigationExitW) {
		ui.FillRect(screen,
			&ui.FillRectOpts{
				Rect:  image.Rect(x-dh, y+h/2-dw, x+dh, y+h/2+dw),
				Color: doorColor,
			},
		)
	}
}

var _ ui.Drawer = (*MapRoomUI)(nil)

// Box default single color debug box
type Box struct {
	Color color.Color
}

var _ ui.Drawer = (*Box)(nil)

func (b *Box) Draw(screen *ebiten.Image, frame image.Rectangle, view *ui.View) {
	ui.FillRect(screen, &ui.FillRectOpts{
		Rect:  frame,
		Color: b.Color,
	})
}
