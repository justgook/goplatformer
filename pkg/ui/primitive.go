package ui

import (
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	g    graphic
	once sync.Once
)

type graphic struct {
	imgOfAPixel *ebiten.Image
}

func (g *graphic) setup() {
	once.Do(func() {
		g.imgOfAPixel = ebiten.NewImage(1, 1)
	})
}

type FillRectOpts struct {
	Rect  image.Rectangle
	Color color.Color
}

func FillRect(target *ebiten.Image, opts *FillRectOpts) {
	g.setup()
	r, c := &opts.Rect, &opts.Color
	g.imgOfAPixel.Fill(*c)
	op := &ebiten.DrawImageOptions{}
	w, h := r.Size().X, r.Size().Y
	op.GeoM.Translate(float64(r.Min.X)*(1/float64(w)), float64(r.Min.Y)*(1/float64(h)))
	op.GeoM.Scale(float64(w), float64(h))
	target.DrawImage(g.imgOfAPixel, op)
}

type DrawRectOpts struct {
	Rect        image.Rectangle
	Color       color.Color
	StrokeWidth int
}

func DrawRect(target *ebiten.Image, opts *DrawRectOpts) {
	g.setup()
	r, c, sw := &opts.Rect, &opts.Color, opts.StrokeWidth
	FillRect(target, &FillRectOpts{
		Rect: image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+sw), Color: *c,
	})

	FillRect(target, &FillRectOpts{
		Rect:  image.Rect(r.Max.X-sw, r.Min.Y, r.Max.X, r.Max.Y),
		Color: *c,
	})

	FillRect(target, &FillRectOpts{
		Rect: image.Rect(r.Min.X, r.Max.Y-sw, r.Max.X, r.Max.Y), Color: *c,
	})

	FillRect(target, &FillRectOpts{
		Rect: image.Rect(r.Min.X, r.Min.Y, r.Min.X+sw, r.Max.Y), Color: *c,
	})
}

// DrawNinePatches based on https://github.com/hajimehoshi/ebiten/blob/main/examples/ui/main.go
func DrawNinePatches(src *ebiten.Image, dst *ebiten.Image, srcRect image.Rectangle, dstRect image.Rectangle) {
	srcX := srcRect.Min.X
	srcY := srcRect.Min.Y
	srcW := srcRect.Dx()
	srcH := srcRect.Dy()

	dstX := dstRect.Min.X
	dstY := dstRect.Min.Y
	dstW := dstRect.Dx()
	dstH := dstRect.Dy()

	op := &ebiten.DrawImageOptions{}
	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			op.GeoM.Reset()

			sx := srcX
			sy := srcY
			sw := srcW / 4
			sh := srcH / 4
			dx := 0
			dy := 0
			dw := sw
			dh := sh
			switch i {
			case 1:
				sx = srcX + srcW/4
				sw = srcW / 2
				dx = srcW / 4
				dw = dstW - 2*srcW/4
			case 2:
				sx = srcX + 3*srcW/4
				dx = dstW - srcW/4
			}
			switch j {
			case 1:
				sy = srcY + srcH/4
				sh = srcH / 2
				dy = srcH / 4
				dh = dstH - 2*srcH/4
			case 2:
				sy = srcY + 3*srcH/4
				dy = dstH - srcH/4
			}

			op.GeoM.Scale(float64(dw)/float64(sw), float64(dh)/float64(sh))
			op.GeoM.Translate(float64(dx), float64(dy))
			op.GeoM.Translate(float64(dstX), float64(dstY))
			dst.DrawImage(src.SubImage(image.Rect(sx, sy, sx+sw, sy+sh)).(*ebiten.Image), op)
		}
	}
}
