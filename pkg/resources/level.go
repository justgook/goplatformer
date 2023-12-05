package resources

import (
	"encoding"
	"image"
	"log/slog"

	. "github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/util"
)

type Level struct {
	*LevelData
	Tilesets []Tileset
}

type Tileset struct {
	GridSize uint
	Image    image.Image
}
type tilesetMarshalBinary struct {
	GridSize uint
	Image    []byte
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (t *Tileset) UnmarshalBinary(input []byte) error {
	target := &tilesetMarshalBinary{}
	if err := Load(input, target); err != nil {
		return util.Catch(err)
	}
	img, err := BytesToImage(target.Image)
	if err != nil {
		return util.Catch(err)
	}
	t.GridSize = target.GridSize
	t.Image = img

	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (t *Tileset) MarshalBinary() (data []byte, err error) {
	img, err := ImageToBytes(t.Image)
	if err != nil {
		return nil, util.Catch(err)
	}
	t.Image = nil
	output, err := Save(tilesetMarshalBinary{
		GridSize: t.GridSize,
		Image:    img,
	})
	if err != nil {
		return nil, util.Catch(err)
	}
	slog.Info("Tileset.MarshalBinary", "GridSize", t.GridSize)

	return output, nil
}

type LevelData struct {
	Rooms        []*Room
	RoomsByExits map[RoomNavigation][]uint
}

type levelMarshalBinary struct {
	*LevelData
	Tilesets []Tileset
}

func (l *Level) UnmarshalBinary(input []byte) error {
	target1 := &levelMarshalBinary{}
	if err := Load(input, target1); err != nil {
		return util.Catch(err)
	}

	l.LevelData = target1.LevelData
	l.Tilesets = target1.Tilesets
	return nil
}

func (l *Level) MarshalBinary() ([]byte, error) {
	output, err := Save(levelMarshalBinary{
		LevelData: l.LevelData,
		Tilesets:  l.Tilesets,
	})
	if err != nil {
		return nil, util.Catch(err)
	}

	return output, nil
}

type Tile struct {
	X int64
	Y int64
	T int64
}

type Room struct {
	Layers            [][]Tile
	RoomNavigation    RoomNavigation
	W                 int
	H                 int
	Collision         []*Object
	TriggerSpawnEnemy []*TriggerSpawnEnemy
	LevelEnter        *LevelEnter
}
type TriggerSpawnEnemy struct {
	Area    image.Rectangle
	Enemies []*Enemy
}

type Enemy struct {
	image.Point
	Patrol []image.Point
}
type LevelEnter struct {
	Start  image.Point
	EnterN image.Point
	EnterE image.Point
	EnterS image.Point
	EnterW image.Point
}

var _ encoding.BinaryMarshaler = (*Level)(nil)
var _ encoding.BinaryUnmarshaler = (*Level)(nil)
var _ encoding.BinaryMarshaler = (*Tileset)(nil)
var _ encoding.BinaryUnmarshaler = (*Tileset)(nil)
