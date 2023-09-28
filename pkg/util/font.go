package util

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func LoadFont(data []byte, size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

