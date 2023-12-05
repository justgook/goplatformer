package level

import (
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/gen"
	"github.com/justgook/goplatformer/pkg/pcg"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/justgook/goplatformer/pkg/resources"
	"github.com/justgook/goplatformer/pkg/util"
)

type MinimapData struct {
	RoomPoint
	Exits util.Bits
	Goal  bool
	Start bool
}

type Data struct {
	Rooms         map[image.Point]*PlayData
	CurrentRoomXY image.Point
	CurrentRoom   *PlayData
	Size          image.Point
}

type RoomPoint = image.Point

type RoomRender struct {
	Render []*ebiten.Image
	Rect   image.Rectangle
}

type TmpRoomResult struct {
	*RoomRender
	Collision         []*Object
	TriggerSpawnEnemy []string
}

type TileData struct {
	Image *ebiten.Image
	Rect  image.Rectangle
}

func (l *Data) fillPlayRoomData(
	rnd *pcg.PCG32,
	input *resources.Level,
	tiles Tiles,
) gen.RoomsLayoutFN {
	return func(info *gen.RoomMetaData, _ image.Point, _ image.Point) {
		available, ok := input.RoomsByExits[info.RoomNavigation]
		if !ok {
			available = input.RoomsByExits[info.RoomNavigation&0x0F]
		}
		randomI := rnd.Bounded(uint32(len(available)))
		resData := input.Rooms[available[randomI]]

		l.Rooms[info.RoomPoint] = &PlayData{
			Features:          info.RoomNavigation,
			RenderLayer:       l.bakeRoomRender(tiles, resData),
			Collision:         util.PointerSliceClone(resData.Collision),
			TriggerSpawnEnemy: util.PointerSliceClone(resData.TriggerSpawnEnemy),
			LevelEnter:        &resources.LevelEnter{},
		}

		*l.Rooms[info.RoomPoint].LevelEnter = *resData.LevelEnter

		// TODO Update to real math
		l.Size.X = max(l.Size.X, resData.W)
		l.Size.Y = max(l.Size.Y, resData.H)
	}
}

func (l *Data) SetCurrent(p image.Point, space *ObjectSpace) {
	space.UnregisterAllObjects()
	l.CurrentRoomXY = p
	l.CurrentRoom = l.Rooms[p]
	slog.Info("components.level.SetCurrent", "collision", l.CurrentRoom.Collision)

	for _, obj := range l.CurrentRoom.Collision {
		space.Add(obj)
	}

	for i := range l.CurrentRoom.TriggerSpawnEnemy {
		obj := l.CurrentRoom.TriggerSpawnEnemy[i].Area
		space.Add(
			resolv.NewObject(
				float64(obj.Min.X),
				float64(obj.Min.Y),
				float64(obj.Dx()),
				float64(obj.Dy()),
				ObjectTagEnemyTrigger,
			))
	}
}

func New(
	rnd *pcg.PCG32,
	goalDistance, branchLength int,
	data *resources.Level,
	collisionSpace *ObjectSpace,
) *Data {
	output := &Data{}
	output.Rooms = make(map[image.Point]*PlayData)
	tiles := remapTileset(data.Tilesets)
	gen.RoomsLayout(rnd, goalDistance, branchLength, output.fillPlayRoomData(rnd, data, tiles))

	collisionSpace.Resize(output.Size.X, output.Size.Y)
	for k := range output.Rooms {
		if output.Rooms[k].Features.Has(RoomNavigationStart) {
			output.SetCurrent(k, collisionSpace)
			break
		}
	}

	return output
}

func (l *Data) bakeRoomRender(tiles Tiles, source *resources.Room) []*ebiten.Image {
	roomLayers := make([]*ebiten.Image, 0, len(source.Layers))

	for _, layer := range source.Layers {
		layerImage := ebiten.NewImage(source.W, source.H)
		for _, tileData := range layer {
			tile := tiles[tileData.T]
			tileImg := tile.Image.SubImage(tile.Rect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tileData.X), float64(tileData.Y))
			layerImage.DrawImage(tileImg, op)
		}
		roomLayers = append(roomLayers, layerImage)
	}
	return roomLayers
}

type PlayData struct {
	Features    RoomNavigation
	RenderLayer []*ebiten.Image
	// TODO merge those two together?
	Collision         []*Object
	TriggerSpawnEnemy []*resources.TriggerSpawnEnemy
	LevelEnter        *resources.LevelEnter
}
