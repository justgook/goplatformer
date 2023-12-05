package main

import (
	"flag"
	"fmt"
	"image"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	. "github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/gameLogger/cli"
	"github.com/justgook/goplatformer/pkg/ldtk/v2"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/justgook/goplatformer/pkg/resources"
	"github.com/justgook/goplatformer/pkg/util"
)

func main() {
	handler := cli.New(os.Stderr, &cli.Options{
		HandlerOptions: slog.HandlerOptions{Level: slog.LevelDebug},
	})
	slog.SetDefault(slog.New(handler))

	inputLevel := flag.String("level", "unknown.ldtk", "img data file")
	outPath := flag.String("o", "unknown.gob", "output location")
	flag.Parse()

	dataBytes := util.GetOrDie(os.ReadFile(*inputLevel))
	jsonData := util.GetOrDie(ldtk.UnmarshalLdtkJSON(dataBytes))

	output := &resources.Level{
		LevelData: &resources.LevelData{Rooms: []*resources.Room{}, RoomsByExits: map[RoomNavigation][]uint{}},
		Tilesets:  []resources.Tileset{},
	}

	slog.Info("========================================================================================")
	tilesetDir := filepath.Dir(*outPath)
	for _, t := range jsonData.Defs.Tilesets {
		if t.RelPath == nil {
			continue
		}
		filename := path.Join(tilesetDir, *t.RelPath)
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))
		filename += ".png"
		imgBytes := util.GetOrDie(os.ReadFile(filename))

		output.Tilesets = append(output.Tilesets, resources.Tileset{
			GridSize: uint(t.TileGridSize),
			Image:    util.GetOrDie(resources.BytesToImage(imgBytes)),
		})
	}

	for worldI, world := range jsonData.Worlds {
		if worldI > 0 {
			panic("game can have only one world (for now)")
		}
		worldLogger := slog.With(
			slog.Group("world",
				slog.String("name", world.Identifier),
			),
		)

		for index, room := range world.Levels {
			//slog.With("room", room.Identifier)
			collisionFound := false
			outputRoom := &resources.Room{
				Layers:            [][]resources.Tile{},
				RoomNavigation:    0,
				W:                 int(room.PxWid),
				H:                 int(room.PxHei),
				Collision:         nil,
				TriggerSpawnEnemy: []*resources.TriggerSpawnEnemy{},
				LevelEnter:        &resources.LevelEnter{},
			}
			roomLogger := worldLogger.With(
				slog.Group("room",
					"name", room.Identifier,
					"attrs", room.FieldInstances,
				),
			)
			//the 1st layer is the top-most and the last is behind.
			slices.Reverse(room.LayerInstances)
			spawnCache := SpawnTmp{}
			for _, layer := range room.LayerInstances {
				var outputLayer []resources.Tile
				layerLogger := roomLogger.With(slog.Group("layer",
					"name", layer.Identifier,
					"GridSize", layer.GridSize,
					"IntGrid", len(layer.IntGridCSV),
					"Entities", len(layer.EntityInstances),
					"AutoLayerTiles", len(layer.AutoLayerTiles),
					"GridTiles", len(layer.GridTiles),
				))

				if len(layer.IntGridCSV) > 0 {
					layerLogger.Info("parsing Collision")
					if collisionFound {
						panic("room can have only one collision room (for now)")
					}
					parseCollisionLayer(layer.IntGridCSV, layer.CWid, layer.CHei, layer.GridSize, outputRoom)
					collisionFound = true
				}

				if len(layer.EntityInstances) > 0 {
					layerLogger.Info("parsing Entities")
					for _, entity := range layer.EntityInstances {
						parseEntity(spawnCache, entity, layer.GridSize, outputRoom)
					}
				}

				if len(layer.AutoLayerTiles) > 0 {
					layerLogger.Info("parsing AutoLayerTiles")
					for _, tile := range layer.AutoLayerTiles {
						outputLayer = append(outputLayer, resources.Tile{
							X: tile.Px[0],
							Y: tile.Px[1],
							T: tile.T,
						})
					}
				}

				if len(layer.GridTiles) > 0 {
					layerLogger.Info("parsing GridTiles")
					for _, tile := range layer.GridTiles {
						outputLayer = append(outputLayer, resources.Tile{
							X: tile.Px[0],
							Y: tile.Px[1],
							T: tile.T,
						})
					}
				}

				if len(outputLayer) > 0 {
					outputRoom.Layers = append(outputRoom.Layers, outputLayer)
				}
			}
			output.Rooms = append(output.Rooms, outputRoom)

			output.RoomsByExits[outputRoom.RoomNavigation] = append(output.RoomsByExits[outputRoom.RoomNavigation], uint(index))
		}
	}
	slog.Info("========================================================================================")
	toFile := util.GetOrDie(resources.Save(output))
	file, _ := os.Create(*outPath)
	defer file.Close()
	util.GetOrDie(file.Write(toFile))
}

/* =======================Parse Entity layer================================================================= */

func parseEntity(spawnCache SpawnTmp, entity ldtk.EntityInstance, cellSize int64, target *resources.Room) {
	switch entity.Identifier {
	case "ExitNorth":
		target.RoomNavigation = target.RoomNavigation.Set(RoomNavigationExitN)

		obj := resolv.NewObject[ObjectTag](float64(entity.Px[0]), float64(entity.Px[1]), float64(entity.Width), float64(entity.Height), ObjectTagExitTriggerNorth)
		target.Collision = append(target.Collision, obj)
	case "ExitEast":
		target.RoomNavigation = target.RoomNavigation.Set(RoomNavigationExitE)

		obj := resolv.NewObject[ObjectTag](float64(entity.Px[0]), float64(entity.Px[1]), float64(entity.Width),
			float64(entity.Height), ObjectTagExitTriggerEast)
		target.Collision = append(target.Collision, obj)
	case "ExitSouth":
		target.RoomNavigation = target.RoomNavigation.Set(RoomNavigationExitS)

		obj := resolv.NewObject[ObjectTag](float64(entity.Px[0]), float64(entity.Px[1]), float64(entity.Width),
			float64(entity.Height), ObjectTagExitTriggerSouth)
		target.Collision = append(target.Collision, obj)
	case "ExitWest":
		target.RoomNavigation = target.RoomNavigation.Set(RoomNavigationExitW)

		obj := resolv.NewObject[ObjectTag](float64(entity.Px[0]), float64(entity.Px[1]), float64(entity.Width),
			float64(entity.Height), ObjectTagExitTriggerWest)
		target.Collision = append(target.Collision, obj)

	case "LevelStart":
		target.RoomNavigation = target.RoomNavigation.Set(RoomNavigationStart)
		target.LevelEnter.Start = EntityToBottomCenterPoint(entity)
	case "LevelGoal":
		target.RoomNavigation = target.RoomNavigation.Set(RoomNavigationGoal)

		obj := resolv.NewObject[ObjectTag](
			float64(entity.Px[0]),
			float64(entity.Px[1]),
			float64(entity.Width),
			float64(entity.Height), ObjectTagGoal)
		target.Collision = append(target.Collision, obj)

	case "PlayerSpawnNorth":
		target.LevelEnter.EnterN = EntityToBottomCenterPoint(entity)
	case "PlayerSpawnEast":
		target.LevelEnter.EnterE = EntityToBottomCenterPoint(entity)
	case "PlayerSpawnSouth":
		target.LevelEnter.EnterS = EntityToBottomCenterPoint(entity)
	case "PlayerSpawnWest":
		target.LevelEnter.EnterW = EntityToBottomCenterPoint(entity)
	case "TriggerSpawnEnemy":
		parseTriggerSpawnEnemy(spawnCache, entity, target)
	case "SpawnEnemyGround":
		parseSpawnEnemyGround(spawnCache, entity, cellSize, target)
	}
}

func EntityToBottomCenterPoint(entity ldtk.EntityInstance) image.Point {
	return image.Point{
		X: int(entity.Px[0]) + int(entity.Width/2),
		Y: int(entity.Px[1]) + int(entity.Height),
	}
}

type SpawnTmp = map[string]*resources.TriggerSpawnEnemy

func getFromSpawnTmp(id string, spawnCache SpawnTmp, target *resources.Room) *resources.TriggerSpawnEnemy {
	output, ok := spawnCache[id]
	if !ok {
		output = &resources.TriggerSpawnEnemy{
			Area:    image.Rectangle{},
			Enemies: []*resources.Enemy{},
		}
		spawnCache[id] = output
		target.TriggerSpawnEnemy = append(target.TriggerSpawnEnemy, output)
	}

	return output
}
func parseSpawnEnemyGround(spawnCache SpawnTmp, entity ldtk.EntityInstance, cellSize int64, target *resources.Room) {
	result := &resources.Enemy{
		Point: image.Point{
			X: int(entity.Px[0]),
			Y: int(entity.Px[1]),
		},
		Patrol: []image.Point{},
	}
	slog.Info("-------------------------------------------------------")
	fmt.Printf("Patrol: %+v\n", entity)
	slog.Info("-------------------------------------------------------")

	for _, value := range entity.FieldInstances {
		switch value.Identifier {
		case "Patrol":
			props := value.Value.([]interface{})
			result.Patrol = append(result.Patrol, image.Point{
				X: int(entity.Px[0]),
				Y: int(entity.Px[1]),
			})
			for _, point := range props {
				p := point.(map[string]interface{})
				result.Patrol = append(result.Patrol, image.Point{
					X: int(p["cx"].(float64)) * int(cellSize),
					Y: int(p["cy"].(float64)) * int(cellSize),
				})
			}

		case "Trigger":
			props := value.Value.(map[string]interface{})
			container := getFromSpawnTmp(props["entityIid"].(string), spawnCache, target)
			container.Enemies = append(container.Enemies, result)
		}
	}

	// getSpawnCache()

}
func parseTriggerSpawnEnemy(spawnCache SpawnTmp, entity ldtk.EntityInstance, target *resources.Room) {
	x := int(entity.Px[0])
	y := int(entity.Px[1])
	w := int(entity.Width)
	h := int(entity.Height)
	area := image.Rect(x, y, x+w, y+h)
	spawnArea := getFromSpawnTmp(entity.Iid, spawnCache, target)
	spawnArea.Area = area
}

/* =======================Parse Collison layer================================================================= */

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

func parseCollisionLayer(input []int64, w, h, cellSize int64, target *resources.Room) {
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
			neighbor.IsSameTags([]ObjectTag{tag}) {
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
			neighbor.IsSameTags([]ObjectTag{tag}) {
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
			slog.Info("parseCollisionLayer:out of bounds(X)",
				"size", len(tmp),
				"X", x,
				"Y", y,
			)
			continue
		}
		if y >= int64(len(tmp[x])) {
			slog.Info("parseCollisionLayer:out of bounds(Y)",
				"size", len(tmp[x]),
				"X", x,
				"Y", y,
			)
			continue
		}

		obj := resolv.NewObject[ObjectTag](
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

	added := make(map[*Object]bool)

	for i := range tmp {
		for _, b := range tmp[i] {
			if b == nil || added[b] {
				continue
			}

			target.Collision = append(target.Collision, b)
			added[b] = true
		}
	}
}
