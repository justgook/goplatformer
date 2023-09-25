package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/justgook/goplatformer/pkg/aseprite"
	"github.com/justgook/goplatformer/pkg/game"
)

func main() {
	outPath := flag.String("o", "unknown.gob", "output location")
	dataPath := flag.String("data", "unknown.json", "json data file")
	imgPath := flag.String("sprite", "unknown.png", "img data file")
	flag.Parse()

	dataBytes, err := os.ReadFile(*dataPath)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	data, err := aseprite.UnmarshalAnimation(dataBytes)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	imageBytes, err := os.ReadFile(*imgPath)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	file, _ := os.Create(*outPath)
	defer file.Close()

	encoder := gob.NewEncoder(file)
	encoder.Encode(game.SpritesheetRaw{
		Image:    imageBytes,
		AnimData: convertAsprite2AnimDataMap(&data),
	})
}

// setup aseprite to:
// Item File name `{layer} {frame}`
// Item Tag Name `{tag}`
func convertAsprite2AnimDataMap(input *aseprite.Animation) game.AnimDataMap {
	output := make(game.AnimDataMap, len(input.Meta.FrameTags))
	layerCount := len(input.Meta.Layers)
	frameMaps := input.Frames
	for _, anim := range input.Meta.FrameTags {
		output[anim.Name] = make([]game.SpriteRawFrame, 0, anim.To-anim.From+1)
		for i := anim.From; i <= anim.To; i++ {
			item := game.SpriteRawFrame{Layers: make([]game.FrameDrawData, 0, layerCount)}
			for _, layer := range input.Meta.Layers {
				key := fmt.Sprintf("%s %d", layer.Name, i)
				data := game.FrameDrawData{
					W:  frameMaps[key].Frame.W,
					H:  frameMaps[key].Frame.H,
					X0: frameMaps[key].Frame.X,
					Y0: frameMaps[key].Frame.Y,
					X1: frameMaps[key].Frame.X + frameMaps[key].Frame.W,
					Y1: frameMaps[key].Frame.Y + frameMaps[key].Frame.H,
					TX: frameMaps[key].SpriteSourceSize.X,
					TY: frameMaps[key].SpriteSourceSize.Y,
				}
				item.Duration = int(frameMaps[key].Duration)
				item.Layers = append(item.Layers, data)
			}
			output[anim.Name] = append(output[anim.Name], item)
		}
	}

	return output
}
