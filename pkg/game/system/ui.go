package system

import (
	"fmt"
	// "image/color"

	"github.com/ebitenui/ebitenui"
	// mock "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer"
	"github.com/justgook/goplatformer/pkg/game/system/ui"
	"github.com/justgook/goplatformer/pkg/resources"

	"image"
)

type UI struct {
	res        *resources.UI
	root       *ebitenui.UI
	StatsPanel *ui.StatsPanel
}

func (u *UI) Init() {
	u.res = goplatformer.UIResources
	u.root = &ebitenui.UI{
		Container: u.res.Root(),
		DisableDefaultFocus: false,
	}
	/* Init Panels */

	// rw := u.root.AddWindow(modalWindow(u.res, ui.NewMap(u.res)))
	//u.root.IsWindowOpen()
	// _ = rw
	//u.root.Container.AddChild(levelMap)

}

func (u *UI) Terminate() {
}

func (u *UI) Update() {
	u.root.Update()
}

func (u *UI) Draw(screen *ebiten.Image) {
	u.root.Draw(screen)
}

func modalWindow(res *resources.UI, c *widget.Container) *widget.Window {
	titleBar := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.TitleBar),
		widget.ContainerOpts.Layout(widget.NewGridLayout(widget.GridLayoutOpts.Columns(3), widget.GridLayoutOpts.Stretch([]bool{true, false, false}, []bool{true}), widget.GridLayoutOpts.Padding(widget.Insets{
			Left:   30,
			Right:  5,
			Top:    6,
			Bottom: 5,
		}))))

	titleBar.AddChild(widget.NewText(
		widget.TextOpts.Text("Modal Window", res.Text.TitleFace, res.Text.IdleColor),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))
	window := widget.NewWindow(
		widget.WindowOpts.Modal(),
		//Set how to close the window. CLICK_OUT will close the window when clicking anywhere
		//that is not a part of the window object
		widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.TitleBar(titleBar, 30),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(200, 200),
		widget.WindowOpts.MaxSize(400, 400),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Resize: ", args.Rect)
		}),
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Move: ", args.Rect)
		}),
	)
	w, h := window.Contents.PreferredSize()
	//Create a rect with the preferred size of the content
	r := image.Rect(0, 0, w, h)
	//Use the Add method to move the window to the specified point
	r = r.Add(image.Point{X: 100, Y: 50})
	//Set the windows location to the rect.
	window.SetLocation(r)
	return window
}
