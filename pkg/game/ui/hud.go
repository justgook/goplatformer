package ui

import (
	"image/color"

	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/font"
)

func NewHUD(face font.Face) *widget.Container {
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout(widget.StackedLayoutOpts.Padding(widget.NewInsetsSimple(25)))),
	)

	// innerContainer := widget.NewContainer(
	// 	widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 0, 0, 255})),
	// 	widget.ContainerOpts.WidgetOpts(
	// 		widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
	// 			HorizontalPosition: widget.AnchorLayoutPositionCenter,
	// 			VerticalPosition:   widget.AnchorLayoutPositionCenter,
	// 			StretchHorizontal:  true,
	// 			StretchVertical:    false,
	// 		}),
	// 		widget.WidgetOpts.MinSize(100, 100),
	// 	),
	// )
	rootContainer.AddChild(Health(face))
	rootContainer.AddChild(Health2(face))
	// _ = innerContainer
	// rootContainer.AddChild(innerContainer)

	return rootContainer
}

func Wrap(child widget.PreferredSizeLocateableWidget) widget.PreferredSizeLocateableWidget {
	btnContainer := widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	btnContainer.AddChild(child)

	return btnContainer
}

func Health(face font.Face) widget.PreferredSizeLocateableWidget {
	label1 := widget.NewText(
		widget.TextOpts.Text("text", face, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)

	return Wrap(label1)
}
func Health2(face font.Face) widget.PreferredSizeLocateableWidget {
	label1 := widget.NewText(
		widget.TextOpts.Text("100/100", face, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)
	// delme := new(string)
	// go func() {
	// 	for {
	// 		time.Sleep(time.Second)
	// 		*delme = fmt.Sprintf("FPS: %f", ebiten.ActualFPS())
	// 	}
	// }()
	// delme = &label1.Label

	return Wrap(label1)
}
