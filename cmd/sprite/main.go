package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/justgook/goplatformer/pkg/aseprite"
	"github.com/justgook/goplatformer/pkg/bin"
	"github.com/justgook/goplatformer/pkg/gameLogger/cli"
	"github.com/justgook/goplatformer/pkg/util"
)

func main() {
	handler := cli.New(os.Stderr, &cli.Options{
		HandlerOptions: slog.HandlerOptions{Level: slog.LevelDebug},
	})
	slog.SetDefault(slog.New(handler))
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

	output := bin.AnimatedSprite{
		Image: imageBytes,
		Data:  convertAseprite2AnimDataMap(&data),
	}

	for _, frame := range data.Frames {
		output.W = int(frame.SourceSize.W)
		output.H = int(frame.SourceSize.H)

		break
	}

	toFile := util.GetOrDie(output.Save())
	file, _ := os.Create(*outPath)
	defer file.Close()
	util.GetOrDie(file.Write(toFile))
}

// setup aseprite to:
// Item File name `{layer} {frame}`
// Item Tag Name `{tag}`
func convertAseprite2AnimDataMap(input *aseprite.Animation) bin.AnimatedSpriteDataMap {
	output := make(bin.AnimatedSpriteDataMap, len(input.Meta.FrameTags))
	layerCount := len(input.Meta.Layers)
	frameMaps := input.Frames
	for _, anim := range input.Meta.FrameTags {
		output[anim.Name] = make([]bin.AnimatedSpriteFrame, 0, anim.To-anim.From+1)
		for i := anim.From; i <= anim.To; i++ {
			item := bin.AnimatedSpriteFrame{Layers: make([]bin.AnimatedSpriteFrameLayer, 0, layerCount)}
			for _, layer := range input.Meta.Layers {
				key := fmt.Sprintf("%s %d", layer.Name, i)
				data := bin.AnimatedSpriteFrameLayer{
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
