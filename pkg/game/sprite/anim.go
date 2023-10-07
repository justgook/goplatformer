package sprite

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	bin "github.com/justgook/goplatformer/pkg/resources"
	"github.com/justgook/goplatformer/pkg/util"
	"image"
	_ "image/png"
)

type Animated struct {
	Sprite       *ebiten.Image
	source       *ebiten.Image
	data         bin.AnimatedSpriteDataMap
	FlipH        bool
	beforeFlip   *ebiten.Image
	width        int
	height       int
	name         string
	currentFrame int
}

func (a *Animated) SetName(name string) {
	a.name = name
}
func (a *Animated) Update() {
	animationFrame := (a.currentFrame / 5) % (len(a.data[a.name]))
	currentFrameData := a.data[a.name][animationFrame]
	a.beforeFlip.Clear()

	op := &ebiten.DrawImageOptions{}
	for _, info := range currentFrameData.Layers {
		rect := image.Rect(info.X0, info.Y0, info.X1, info.Y1)
		result := a.source.SubImage(rect).(*ebiten.Image)
		//op.GeoM.Translate(float64(info.TX), float64(info.TY))
		op.GeoM.SetElement(0, 2, float64(info.TX))
		op.GeoM.SetElement(1, 2, float64(info.TY))
		a.beforeFlip.DrawImage(result, op)
	}

	if a.FlipH {
		op.GeoM.Reset()
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(a.width), 0)
	}
	a.Sprite.Clear()
	a.Sprite.DrawImage(a.beforeFlip, op)
	a.currentFrame++
}

func (a *Animated) FromRaw(input *bin.AnimatedSprite) {
	img, _ := util.Get2OrDie(image.Decode(bytes.NewReader(input.Image)))
	a.source = ebiten.NewImageFromImage(img)
	a.data = input.Data
	a.width = input.W
	a.height = input.H
	a.Sprite = ebiten.NewImage(input.W, input.H)
	a.beforeFlip = ebiten.NewImage(input.W, input.H)

}

func (a *Animated) Load(input []byte) error {
	target := &bin.AnimatedSprite{}
	if err := bin.Load(input, target); err != nil {
		return util.Catch(err)
	}
	a.FromRaw(target)

	return nil
}
