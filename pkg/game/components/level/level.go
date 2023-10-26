package level

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/pcg"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/justgook/goplatformer/pkg/resources"
	"github.com/justgook/goplatformer/pkg/util"
)

type Data struct {
	RoomsInfo    []RoomMetaData
	RoomsResults []*RoomResult
	Size         image.Point
	Goal         image.Point
	Start        image.Point
}
type RoomPoint = image.Point
type RoomMetaData struct {
	Exits util.Bits
	RoomPoint
	Goal  bool
	Start bool
}

type TagType = resources.CollisionTag
type RoomResult struct {
	Render    []*ebiten.Image
	Rect      image.Rectangle
	Collision []*resolv.Object[TagType]
}

const (
	ExitN util.Bits = 1 << iota
	ExitE
	ExitS
	ExitW
)

func New(rnd *pcg.PCG32, goalDistance, branchLength int, data *resources.Level, space *resolv.Space[string]) *Data {
	output := generateData(rnd, goalDistance, branchLength)
	tileset := ebiten.NewImageFromImage(data.Image)
	output.RoomsResults = make([]*RoomResult, len(output.RoomsInfo))
	output.Size = image.Point{}

	for i := range output.RoomsInfo {
		info := output.RoomsInfo[i]
		available := data.RoomsByExits[info.Exits]
		randomI := rnd.Bounded(uint32(len(available)))
		resData := data.Rooms[available[randomI]]
		room := &RoomResult{
			Collision: util.PointerSliceClone(resData.Collision),
			Render:    bakeRoomLayers(tileset, resData),
			Rect: image.Rect(
				info.RoomPoint.X*256,           // TODO update to real coordinate
				info.RoomPoint.Y*256,           // TODO update to real coordinate
				info.RoomPoint.X*256+resData.W, // TODO update to real coordinate
				info.RoomPoint.Y*256+resData.H, // TODO update to real coordinate
			),
		}

		output.RoomsResults[i] = room

		output.Size.X = max(output.Size.X, room.Rect.Max.X)
		output.Size.Y = max(output.Size.Y, room.Rect.Max.Y)
	}

	space.Resize(output.Size.X, output.Size.Y)

	for i := range output.RoomsResults {
		room := output.RoomsResults[i]
		for _, obj := range room.Collision {
			obj.X += float64(room.Rect.Min.X)
			obj.Y += float64(room.Rect.Min.Y)

			addToSpace(space, obj)
		}
	}

	return output
}

func addToSpace(space *resolv.Space[string], obj *resolv.Object[TagType]) {
	tag := "solid"
	switch obj.Tags[0] {
	case resources.CollisionTagSolid:
		tag = "solid"
	case resources.CollisionTagOneWayUp:
		tag = "platform"
	case resources.CollisionTagLava:
		tag = "lava"
	case resources.CollisionTagExitNorth:
		tag = "ExitNorth"
	case resources.CollisionTagExitEast:
		tag = "ExitEast"
	case resources.CollisionTagExitSouth:
		tag = "ExitSouth"
	case resources.CollisionTagExitWest:
		tag = "ExitWest"
	}
	out := resolv.NewObject(obj.X, obj.Y, obj.W, obj.H, tag)
	space.Add(out)
}

func bakeRoomLayers(tileSet *ebiten.Image, source *resources.Room) []*ebiten.Image {
	roomLayers := make([]*ebiten.Image, 0, len(source.Layers))
	for _, layer := range source.Layers {
		layerImage := ebiten.NewImage(source.W, source.H)
		for _, tileData := range layer {
			id := tileData.T
			x1 := int(id%12) * 16
			y1 := int(id/12) * 16
			rect := image.Rect(x1, y1, x1+16, y1+16)
			tileImg := tileSet.SubImage(rect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tileData.X), float64(tileData.Y))
			layerImage.DrawImage(tileImg, op)
		}
		roomLayers = append(roomLayers, layerImage)
	}
	return roomLayers
}
