package ui

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/justgook/goplatformer/pkg/resources"
)

type Map struct {
}

func NewMap(res *resources.UI) *widget.Container {
	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.Image),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(100, 100),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(res.Panel.Padding),
			widget.RowLayoutOpts.Spacing(15))),
	)
	titleText := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text("IM MAP", res.Text.TitleFace, res.Text.IdleColor))
	c.AddChild(titleText)

	return c
}
