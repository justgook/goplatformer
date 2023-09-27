package sprite

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/bin"
	"github.com/justgook/goplatformer/pkg/util"
	"image"
	"image/color"
	_ "image/png"
)

type Animated struct {
	Sprite       *ebiten.Image
	source       *ebiten.Image
	data         bin.AnimatedSpriteDataMap
	name         string
	currentFrame int
}

func (a *Animated) SetName(name string) {
	a.name = name
}
func (a *Animated) Update() {
	animationFrame := (a.currentFrame / 5) % (len(a.data[a.name]))
	currentFrameData := a.data[a.name][animationFrame]
	a.Sprite.Fill(color.RGBA{})
	for _, info := range currentFrameData.Layers {
		rect := image.Rect(info.X0, info.Y0, info.X1, info.Y1)
		result := a.source.SubImage(rect).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(info.TX), float64(info.TY))
		a.Sprite.DrawImage(result, op)
	}
	a.currentFrame++
}

func (a *Animated) FromRaw(input *bin.AnimatedSprite) {
	img, _ := util.Get2OrDie(image.Decode(bytes.NewReader(input.Image)))
	a.data = input.Data
	a.source = ebiten.NewImageFromImage(img)
	a.Sprite = ebiten.NewImage(input.W, input.H)
}

func (a *Animated) Load(data []byte) error {
	delme := &bin.AnimatedSprite{}
	if err := delme.Load(data); err != nil {
		return err
	}
	a.FromRaw(delme)

	return nil
}
