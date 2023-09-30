package resources

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type UI struct {
	Panel *UIPanel
	Text  *UIText
}

type UIPanel struct {
	Image    *image.NineSlice
	TitleBar *image.NineSlice
	Padding  widget.Insets
}

type UIText struct {
	IdleColor     color.Color
	DisabledColor color.Color
	Face          font.Face
	TitleFace     font.Face
	BigTitleFace  font.Face
	SmallFace     font.Face
}

// TODO move it to other place
func (res *UI) Root() *widget.Container {
	root := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewStackedLayout()))

	root.AddChild(res.HUDLayouy())
	// root.AddChild(res.PanelLayouy())
	return root
}
func (res *UI) HUDLayouy() widget.PreferredSizeLocateableWidget {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(8)),
		)),
	)

	topContainer := widget.NewContainer(
		// widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x00, 0xff, 0x00, 0x55})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true, false}, []bool{false}),
		)),
	)

	bottomContainer := widget.NewContainer(
		// widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0xFF, 0x00, 0xff, 0x55})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{false, true, false}, []bool{false, false, false}),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(100, 20),
		),
	)

	topContainer.AddChild(MeterBar())
	topContainer.AddChild(res.NewMinimap())

	bottomContainer.AddChild(res.Mock(48, 32))
	bottomContainer.AddChild(res.BossHealth())
	bottomContainer.AddChild(res.Mock(64, 32))

	rootContainer.AddChild(topContainer)
	rootContainer.AddChild(widget.NewContainer())
	rootContainer.AddChild(bottomContainer)

	return rootContainer
}

func MeterBar() widget.PreferredSizeLocateableWidget {
	barW, barH := 64, 6
	HealthBarW, heakthBarH := 64, 12

	cccc := uint8(0x88)
	trackImages := &widget.ProgressBarImage{
		Idle:  image.NewNineSliceColor(color.NRGBA{cccc, cccc, cccc, 255}),
		Hover: image.NewNineSliceColor(color.NRGBA{cccc, cccc, cccc, 255}),
	}
	rootContainer := widget.NewContainer(
		// widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0xFF, 0x00, 0xff, 0x55})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
		)),
	)

	healtBar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			// Set the minimum size for the progress bar.
			// This is necessary if you wish to have the progress bar be larger than
			// the provided track image. In this exampe since we are using NineSliceColor
			// which is 1px x 1px we must set a minimum size.
			widget.WidgetOpts.MinSize(HealthBarW, heakthBarH ),
		),
		widget.ProgressBarOpts.Images(
			// Set the track images (Idle, Hover, Disabled).
			trackImages,
			// Set the progress images (Idle, Hover, Disabled).
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{220, 86, 75, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		// Set the min, max, and current values.
		widget.ProgressBarOpts.Values(0, 10, 7),
		// Set how much of the track is displayed when the bar is overlayed.
		widget.ProgressBarOpts.TrackPadding(widget.NewInsetsSimple(2)),
	)

	staminaBar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(barW, barH),
		),
		widget.ProgressBarOpts.Images(
			trackImages,
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{107, 252, 185, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		widget.ProgressBarOpts.Values(0, 10, 3),

		widget.ProgressBarOpts.TrackPadding(widget.NewInsetsSimple(1)),

	)
	manaBar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(barW, barH),
		),
		widget.ProgressBarOpts.Images(
			trackImages,
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{86, 176, 239, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		widget.ProgressBarOpts.Values(0, 10, 8),
		widget.ProgressBarOpts.TrackPadding(widget.NewInsetsSimple(1)),

	)
	/*
		To update a progress bar programmatically you can use
		healtBar.SetCurrent(value)
		healtBar.GetCurrent()
		healtBar.Min = 5
		hProgressbar.Max = 10
	*/

	rootContainer.AddChild(healtBar)
	rootContainer.AddChild(staminaBar)
	rootContainer.AddChild(manaBar)
	return rootContainer
}

func (res *UI) Mock(w, h int) widget.PreferredSizeLocateableWidget {
	rootContainer := widget.NewContainer(
		// widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x00, 0x00, 0x66, 0x88})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	c := color.NRGBA{0xcc, 0xcc, 0xcc, 0xff}
	i := ebiten.NewImage(w, h)
	i.Fill(c)
	child := widget.NewGraphic(widget.GraphicOpts.Image(i),
		widget.GraphicOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				// HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition: widget.AnchorLayoutPositionEnd,
			}),
		),
	)
	rootContainer.AddChild(child)
	return rootContainer
}

func (res *UI) NewMinimap() widget.PreferredSizeLocateableWidget {
	rootContainer := widget.NewContainer(
		// widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x00, 0x00, 0x66, 0x88})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	c := color.NRGBA{0xcc, 0xcc, 0xcc, 0xff}
	i := ebiten.NewImage(64, 64)
	i.Fill(c)
	child := widget.NewGraphic(widget.GraphicOpts.Image(i),
		widget.GraphicOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
			}),
		),
	)
	rootContainer.AddChild(child)
	return rootContainer
}
func (res *UI) BossHealth() widget.PreferredSizeLocateableWidget {
	rootContainer := widget.NewContainer(
		// widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x00, 0x00, 0x66, 0x88})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionCenter,
				VerticalPosition:   widget.GridLayoutPositionEnd,
				MaxWidth:           300,
				MaxHeight:          200,
			}),
		),
	)

	healthBar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  true,
			}),
			widget.WidgetOpts.MinSize(20, 24),
		),
		widget.ProgressBarOpts.Images(
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 0x66}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 0x66}),
			},
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{220, 86, 75, 0x66}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 0, 255, 0x66}),
			},
		),
		widget.ProgressBarOpts.Values(0, 10, 8),
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Top:    2,
			Right:  2,
			Bottom: 2,
			Left:   2,
		}),
	)
	rootContainer.AddChild(healthBar)
	return rootContainer
}

func (res *UI) PanelLayouy() widget.PreferredSizeLocateableWidget {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{true}),
		)),
	)

	leftContainer := widget.NewContainer(
		// widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255})),
		widget.ContainerOpts.BackgroundImage(res.Panel.Image),

		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionCenter,
				VerticalPosition:   widget.GridLayoutPositionCenter,
				MaxWidth:           200,
				MaxHeight:          200,
			}),
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	root.AddChild(leftContainer)

	rightContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 255, 255, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	root.AddChild(rightContainer)

	return root
}

