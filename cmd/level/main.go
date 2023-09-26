package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/justgook/goplatformer/pkg/ldtk/v2"
)

func main() {
	outPath := flag.String("o", "unknown.gob", "output location")
	flag.Parse()

	args := flag.Args()
	slog.Info(fmt.Sprintf("BUILDIN from: %v", args))
	dataBytes, err := os.ReadFile(args[0])
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	data, err := ldtk.UnmarshalLdtkJSON(dataBytes)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	slog.Info("all defs", "LevelFields", data.Defs.LevelFields)
	slog.Info("========================================================================================")
	for _, world := range data.Worlds {
		slog.Info("got world:", "world", world.Identifier)
		for _, level := range world.Levels {
			slog.Info("got level:",
				"level", level.Identifier,
				"FieldInstances", level.FieldInstances,
			)
			for _, layer := range level.LayerInstances {
				slog.Info("got layer:",
					"layer", layer.Identifier,
					"IntGrid", len(layer.IntGridCSV),
				)
			}
		}
	}
	slog.Info("========================================================================================")

	file, _ := os.Create(*outPath)
	defer file.Close()

	output := "hello World"
	encoder := gob.NewEncoder(file)
	encoder.Encode(output)
}

