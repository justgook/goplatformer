package util

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadImage(input []byte) (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, Catch(err)
	}
	return ebiten.NewImageFromImage(img), nil
}

