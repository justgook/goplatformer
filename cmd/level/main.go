package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/justgook/goplatformer/pkg/bin"
	"github.com/justgook/goplatformer/pkg/gameLogger/cli"
	"github.com/justgook/goplatformer/pkg/ldtk/v2"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/justgook/goplatformer/pkg/util"
)

type Tag = bin.TagType
type Object = resolv.Object[Tag]

func main() {
	handler := cli.New(os.Stderr, &cli.Options{
		HandlerOptions: slog.HandlerOptions{Level: slog.LevelDebug},
	})
	slog.SetDefault(slog.New(handler))

	outPath := flag.String("o", "unknown.gob", "output location")
	flag.Parse()
	args := flag.Args()

	slog.Info(fmt.Sprintf("BUILDIN from: %v", args))
	dataBytes := util.GetOrDie(os.ReadFile(args[0]))
	jsonData := util.GetOrDie(ldtk.UnmarshalLdtkJSON(dataBytes))

	output := &bin.Level{}
	collisionFound := false
	slog.Info("========================================================================================")
	for worldI, world := range jsonData.Worlds {
		if worldI > 0 {
			panic("game can have only one world (for now)")
		}
		slog.Info("got world:", "world", world.Identifier)
		for _, room := range world.Levels {
			outputRoom := &bin.Room{
				Layers:    [][]int{},
				Doors:     bin.Doors{},
				W:         int(room.PxWid),
				H:         int(room.PxHei),
				Collision: nil,
			}
			slog.Info("got room:",
				"room", room.Identifier,
				"FieldInstances", room.FieldInstances,
			)
			for _, layer := range room.LayerInstances {
				slog.Info("----")
				slog.Info("working with layer",
					"layer", layer.Identifier,
					"IntGrid", len(layer.IntGridCSV),
					"GridSize", layer.GridSize,
				)

				if len(layer.IntGridCSV) > 0 {
					slog.Info("set collision:", "layer", layer.Identifier)
					if collisionFound {
						panic("room can have only one collision room (for now)")
					}
					outputRoom.Collision = intGridToCollision(layer.IntGridCSV, layer.CWid, layer.CHei, layer.GridSize)
					collisionFound = true
				}

				if len(layer.EntityInstances) > 0 {
					slog.Info("set Entities:", "layer", layer.Identifier, "EntityInstances", len(layer.EntityInstances))
				}

				if len(layer.AutoLayerTiles) > 0 {
					slog.Info("set AutoLayerTiles:", "layer", layer.Identifier, "AutoLayerTiles", len(layer.AutoLayerTiles))
				}

				if len(layer.GridTiles) > 0 {
					slog.Info("set GridTiles:", "layer", layer.Identifier, "GridTiles", len(layer.GridTiles))
				}
			}
			output.Rooms = append(output.Rooms, outputRoom)
		}
	}
	slog.Info("========================================================================================")
	toFile := util.GetOrDie(output.Save())
	file, _ := os.Create(*outPath)
	defer file.Close()
	util.GetOrDie(file.Write(toFile))
	// saveFile(*outPath, output)
}

// func saveFile(path string, output any) {
// 	file, _ := os.Create(path)
// 	encoder := gob.NewEncoder(file)
// 	util.OrDie(encoder.Encode(output))
// 	util.OrDie(file.Close())
// }

type Point struct {
	X int64
	Y int64
}

func addToCache(cache map[*Object][]Point, matrix [][]*Object, obj *Object, x, y int64) {
	matrix[x][y] = obj
	if cache[obj] == nil {
		cache[obj] = []Point{}
	}
	cache[obj] = append(cache[obj], Point{X: x, Y: y})
}

func updateCache(cache map[*Object][]Point, matrix [][]*Object, was, now *Object) {
	for _, point := range cache[was] {
		addToCache(cache, matrix, now, point.X, point.Y)
		matrix[point.X][point.Y] = now
	}
	delete(cache, was)
}

func intGridToCollision(input []int64, w, h, cellSize int64) []*Object {
	tmp := make([][]*Object, h)
	for i := range tmp {
		tmp[i] = make([]*Object, w)
	}

	cache := map[*Object][]Point{}

	mergeLeft := func(x, y, tag int64) bool {
		if x <= 0 {
			return false
		}
		current := tmp[x][y]
		neighbor := tmp[x-1][y]
		if neighbor != nil &&
			current.Y == neighbor.Y &&
			current.H == neighbor.H &&
			neighbor.IsSameTags([]Tag{tag}) {
			neighbor.W += float64(cellSize)
			updateCache(cache, tmp, current, neighbor)

			return true
		}

		return false
	}

	mergeUp := func(x, y, tag int64) bool {
		if y <= 0 {
			return false
		}

		neighbor := tmp[x][y-1]
		current := tmp[x][y]

		if neighbor != nil &&
			current.X == neighbor.X &&
			current.W == neighbor.W &&
			neighbor.IsSameTags([]Tag{tag}) {
			neighbor.H += float64(cellSize)

			updateCache(cache, tmp, current, neighbor)

			return true
		}

		return false
	}
	for i, tag := range input {
		if tag == 0 {
			continue
		}

		y := int64(i) / w
		x := int64(i) % w
		obj := resolv.NewObject[Tag](
			float64(x*cellSize),
			float64(y*cellSize),
			float64(cellSize),
			float64(cellSize),
			tag,
		)

		addToCache(cache, tmp, obj, x, y)

		a := mergeUp(x, y, tag)
		b := mergeLeft(x, y, tag)
		if !a && b {
			mergeUp(x, y, tag)

		}
	}

	var output []*Object

	added := make(map[*Object]bool)

	for i := range tmp {
		for _, b := range tmp[i] {
			if b == nil || added[b] {
				continue
			}
			output = append(output, b)
			added[b] = true
		}
	}

	return output
}

