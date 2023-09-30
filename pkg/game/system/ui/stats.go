package ui

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/justgook/goplatformer/pkg/resources"
	"image/color"
)

type StatsPanel struct {
	Res  *resources.UI
	Node widget.PreferredSizeLocateableWidget
}

func (s *StatsPanel) Init() {
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x00, 0xFF, 0xFF, 0xff})),
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			//Which direction to layout children
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			//Set how much padding before displaying content
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			//Set how far apart to space the children
			widget.RowLayoutOpts.Spacing(15),
		)),
	)

	innerContainer1 := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 0, 0, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				//Specify where within the row or column this element should be positioned.
				Position: widget.RowLayoutPositionStart,
				//Should this widget be stretched across the row or column
				Stretch: false,
				//How wide can this element grow to (override preferred widget size)
				MaxWidth: 100,
				//How tall can this element grow to (override preferred widget size)
				MaxHeight: 100,
			}),
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	rootContainer.AddChild(innerContainer1)

	innerContainer2 := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 255, 0, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				//Specify where within the row or column this element should be positioned.
				Position: widget.RowLayoutPositionCenter,
				//Should this widget be stretched across the row or column
				Stretch: true,
				//How wide can this element grow to (override preferred widget size)
				MaxWidth: 200,
				//How tall can this element grow to (override preferred widget size)
				MaxHeight: 100,
			}),
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	rootContainer.AddChild(innerContainer2)

	innerContainer3 := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				//Specify where within the row or column this element should be positioned.
				Position: widget.RowLayoutPositionEnd,
				//Should this widget be stretched across the row or column
				Stretch: true,
				//How wide can this element grow to (override preferred widget size)
				MaxWidth: 400,
				//How tall can this element grow to (override preferred widget size)
				MaxHeight: 100,
			}),
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	rootContainer.AddChild(innerContainer3)
	s.Node = rootContainer
}
