package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

var (
	Debug           = false
	debugColor      = color.RGBA{R: 0xff, A: 0xff}
	debugColorShift = colorm.ColorM{}
)

func debugBorders(screen *ebiten.Image, root containerEmbed) {
	var queue []containerEmbed
	queue = append(queue, root)
	renderColor := resetDebugColor()

	for len(queue) > 0 {
		levelSize := len(queue)
		for levelSize != 0 {
			curr := queue[0]
			queue = queue[1:]

			DrawRect(screen, &DrawRectOpts{
				Rect:        curr.frame,
				Color:       renderColor,
				StrokeWidth: 2,
			})

			for _, c := range curr.children {
				if c.item.Display == DisplayNone {
					continue
				}
				queue = append(queue, c.item.containerEmbed)
			}
			levelSize--
		}

		renderColor = rotateDebugColor()
	}
}

func rotateDebugColor() color.Color {
	debugColorShift.RotateHue(1.66)

	return debugColorShift.Apply(debugColor)
}

func resetDebugColor() color.Color {
	debugColorShift = colorm.ColorM{}
	return debugColor
}
