package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/justgook/goplatformer/pkg/resources"
	"github.com/justgook/goplatformer/pkg/gameLogger/cli"
	"github.com/justgook/goplatformer/pkg/ldtk/v2"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/justgook/goplatformer/pkg/util"
)

type Tag = resources.TagType
type Object = resolv.Object[Tag]

func main() {
	handler := cli.New(os.Stderr, &cli.Options{
		HandlerOptions: slog.HandlerOptions{Level: slog.LevelDebug},
	})
	slog.SetDefault(slog.New(handler))

	tileset := flag.String("tileset", "unknown.png", "img data file")
	inputLevel := flag.String("level", "unknown.ldtk", "img data file")
	outPath := flag.String("o", "unknown.gob", "output location")
	flag.Parse()

	dataBytes := util.GetOrDie(os.ReadFile(*inputLevel))
	jsonData := util.GetOrDie(ldtk.UnmarshalLdtkJSON(dataBytes))

	output := &resources.Level{
		Rooms: []*resources.Room{},
		Image: util.GetOrDie(os.ReadFile(*tileset)),
	}

	slog.Info("========================================================================================")
	for worldI, world := range jsonData.Worlds {
		if worldI > 0 {
			panic("game can have only one world (for now)")
		}
		slog.Info("got world:", "world", world.Identifier)

		for _, room := range world.Levels {
			//slog.With("room", room.Identifier)
			collisionFound := false
			outputRoom := &resources.Room{
				Layers:    [][]resources.Tile{},
				Doors:     resources.Doors{},
				W:         int(room.PxWid),
				H:         int(room.PxHei),
				Collision: nil,
			}

			slog.Info("got room:",
				"room", room.Identifier,
				"FieldInstances", room.FieldInstances,
			)
			slog.Info("------------------------------------------------------")

			for _, layer := range room.LayerInstances {
				outputLayer := []resources.Tile{}
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
					for _, tile := range layer.AutoLayerTiles {
						outputLayer = append(outputLayer, resources.Tile{
							X: tile.Px[0],
							Y: tile.Px[1],
							T: tile.T,
						})
					}
				}

				if len(layer.GridTiles) > 0 {
					slog.Info("set GridTiles:", "layer", layer.Identifier, "GridTiles", len(layer.GridTiles))
				}
				outputRoom.Layers = append(outputRoom.Layers, outputLayer)
				slog.Info("----")
			}
			output.Rooms = append(output.Rooms, outputRoom)
		}
	}
	slog.Info("========================================================================================")
	slog.Info("layers:!!!", "bytes", len(output.Rooms[0].Layers))
	toFile := util.GetOrDie(resources.Save(output))
	file, _ := os.Create(*outPath)
	defer file.Close()
	util.GetOrDie(file.Write(toFile))
}

type Point struct {
	X int64
	Y int64
}

func addToCache(cache map[*Object][]Point, matrix [][]*Object, obj *Object, x, y int64) {
	//slog.Info("addToCache",
	//	"X", x,
	//	"Y", y,
	//)
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
	slog.Info("intGridToCollision",
		"w", w,
		"h", h,
		"cellSize", cellSize,
	)
	tmp := make([][]*Object, w)
	for i := range tmp {
		tmp[i] = make([]*Object, h)
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
		if x >= int64(len(tmp)) {
			slog.Info("intGridToCollision:out of bounds(X)",
				"size", len(tmp),
				"X", x,
				"Y", y,
			)
			continue
		}
		if y >= int64(len(tmp[x])) {
			slog.Info("intGridToCollision:out of bounds(Y)",
				"size", len(tmp[x]),
				"X", x,
				"Y", y,
			)
			continue
		}

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

