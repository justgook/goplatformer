package ui

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/font"
)

func NewMap(face font.Face) *widget.Container {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0x80})),

		widget.ContainerOpts.Layout(widget.NewStackedLayout(widget.StackedLayoutOpts.Padding(widget.NewInsetsSimple(25)))),
	)
	rootContainer.AddChild(Map(face))
	return rootContainer
}

func Map(face font.Face) widget.PreferredSizeLocateableWidget {
	label1 := widget.NewText(
		widget.TextOpts.Text("im map", face, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)

	return Wrap(label1)
}
