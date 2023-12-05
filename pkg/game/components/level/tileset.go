package level

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/resources"
	"image"
)

type Tiles = map[int64]*TileData

func remapTileset(Tilesets []resources.Tileset) Tiles {
	tiles := make(Tiles, len(Tilesets))
	tileId := 0
	for _, v := range Tilesets {
		cols := v.Image.Bounds().Dx() / int(v.GridSize)
		rows := v.Image.Bounds().Dy() / int(v.GridSize)

		img := ebiten.NewImageFromImage(v.Image)
		// TODO move that to the `resources.Level`
		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				s := int(v.GridSize)
				x := col * s
				y := row * s
				tiles[int64(tileId)] = &TileData{
					Image: img,
					Rect:  image.Rect(x, y, x+s, y+s),
				}
				tileId++
			}
		}
	}
	return tiles
}
