package main

import (
	"flag"
	"fmt"
	"image"
	"log/slog"
	"os"
	"time"

	"github.com/justgook/goplatformer/pkg/aseprite"
	"github.com/justgook/goplatformer/pkg/gameLogger/cli"
	"github.com/justgook/goplatformer/pkg/resources"
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

	output := resources.AnimatedSprite{
		Image: imageBytes,
		Anim:  make(map[string]resources.AnimatedSpriteAnimData),
	}

	frameMaps := data.Frames
	for _, tag := range data.Meta.FrameTags {
		frames := make([]*image.Rectangle, 0, tag.To-tag.From)
		durations := make([]time.Duration, 0, tag.To-tag.From)
		for i := tag.From; i <= tag.To; i++ {
			key := fmt.Sprintf("%s %d", tag.Name, i)
			rect := image.Rect(frameMaps[key].Frame.X, frameMaps[key].Frame.Y, frameMaps[key].Frame.X+frameMaps[key].Frame.W, frameMaps[key].Frame.Y+frameMaps[key].Frame.H)
			frames = append(frames, &rect)
			durations = append(durations, time.Millisecond*time.Duration(frameMaps[key].Duration))
		}
		output.Anim[tag.Name] = resources.AnimatedSpriteAnimData{
			Frames:    frames,
			Durations: durations,
		}
	}

	for _, frame := range data.Frames {
		output.W = int(frame.SourceSize.W)
		output.H = int(frame.SourceSize.H)
		break
	}

	toFile := util.GetOrDie(resources.Save(output))
	file, _ := os.Create(*outPath)
	defer file.Close()
	util.GetOrDie(file.Write(toFile))
}
