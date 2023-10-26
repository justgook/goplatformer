package resources

import (
	"encoding"
	"image"

	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/justgook/goplatformer/pkg/util"
)

type Exits = util.Bits

const (
	ExitN util.Bits = 1 << iota
	ExitE
	ExitS
	ExitW
)

type TagType = int64
type Level struct {
	*LevelData
	Image image.Image
}
type LevelData struct {
	Rooms        []*Room
	RoomsByExits map[Exits][]uint
}

type levelMarshalBinary struct {
	LevelData []byte
	Image     []byte
}

func (l *Level) UnmarshalBinary(input []byte) error {
	target1 := &levelMarshalBinary{}
	if err := Load(input, target1); err != nil {
		return util.Catch(err)
	}
	if err := Load(target1.LevelData, &l.LevelData); err != nil {
		return util.Catch(err)
	}
	img, err := BytesToImage(target1.Image)
	if err != nil {
		return util.Catch(err)
	}
	l.Image = img
	return nil
}

func (l *Level) MarshalBinary() ([]byte, error) {
	rooms, err := Save(l.LevelData)
	if err != nil {
		return nil, util.Catch(err)
	}
	img, err := ImageToBytes(l.Image)
	if err != nil {
		return nil, util.Catch(err)
	}
	output, err := Save(levelMarshalBinary{
		LevelData: rooms,
		Image:     img,
	})
	if err != nil {
		return nil, util.Catch(err)
	}

	return output, nil
}

type Doors struct {
	N bool
	E bool
	S bool
	W bool
}

type Tile struct {
	X int64
	Y int64
	T int64
}

type Room struct {
	Layers    [][]Tile
	Exits     Exits
	W         int
	H         int
	Collision []*resolv.Object[TagType]
}

type CollisionTag = int64

const (
	CollisionTagSolid     CollisionTag = 1
	CollisionTagOneWayUp  CollisionTag = 3
	CollisionTagLava      CollisionTag = 4
	CollisionTagExitNorth CollisionTag = 5
	CollisionTagExitEast  CollisionTag = 6
	CollisionTagExitSouth CollisionTag = 7
	CollisionTagExitWest  CollisionTag = 8
)

var _ encoding.BinaryMarshaler = (*Level)(nil)
var _ encoding.BinaryUnmarshaler = (*Level)(nil)
