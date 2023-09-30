package goplatformer

import (
	"bytes"
	_ "embed"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/justgook/goplatformer/pkg/resources"
	baseImg "image"

	"github.com/ebitenui/ebitenui/image"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/util"
	_ "image/png"
)

const (
	backgroundColor = "131a22"

	textIdleColor     = "dff4ff"
	textDisabledColor = "5a7a91"

	labelIdleColor     = textIdleColor
	labelDisabledColor = textDisabledColor

	buttonIdleColor     = textIdleColor
	buttonDisabledColor = labelDisabledColor

	listSelectedBackground         = "4b687a"
	listDisabledSelectedBackground = "2a3944"

	listFocusedBackground = "2a3944"

	headerColor = textIdleColor

	textInputCaretColor         = "e7c34b"
	textInputDisabledCaretColor = "766326"

	toolTipColor = backgroundColor

	separatorColor = listDisabledSelectedBackground
)

var (
	//go:embed asset/ui/res/panel-idle.png
	panelIdle []byte
	//go:embed asset/ui/res/titlebar-idle.png
	titlebar []byte
	//go:embed asset/ui/excel.ttf
	font1       []byte
	UIResources = func() *resources.UI {
		fontFace := util.GetOrDie(util.LoadFont(font1, 20))
		titleFontFace := util.GetOrDie(util.LoadFont(font1, 24))
		bigTitleFontFace := util.GetOrDie(util.LoadFont(font1, 28))
		toolTipFace := util.GetOrDie(util.LoadFont(font1, 15))

		return &resources.UI{
			Panel:     util.GetOrDie(newPanelResources(panelIdle, titlebar)),
			Text: &resources.UIText{
				IdleColor:     util.HexToColor(textIdleColor),
				DisabledColor: util.HexToColor(textDisabledColor),
				Face:          fontFace,
				TitleFace:     titleFontFace,
				BigTitleFace:  bigTitleFontFace,
				SmallFace:     toolTipFace,
			},
		}
	}()
)

func loadImageNineSlice(input []byte, centerWidth int, centerHeight int) (*image.NineSlice, error) {
	img, _, err := baseImg.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, util.Catch(err)
	}
	i := ebiten.NewImageFromImage(img)

	w := i.Bounds().Dx()
	h := i.Bounds().Dy()
	return image.NewNineSlice(i,
			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
		nil
}

func newPanelResources(a, b []byte) (*resources.UIPanel, error) {
	i, err := loadImageNineSlice(a, 10, 10)
	if err != nil {
		return nil, util.Catch(err)
	}
	t, err := loadImageNineSlice(b, 10, 10)
	if err != nil {
		return nil, util.Catch(err)
	}
	return &resources.UIPanel{
		Image:    i,
		TitleBar: t,
		Padding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		},
	}, nil
}
