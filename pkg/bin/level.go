package bin

import (
	"bytes"
	"encoding/gob"
	"image"

	"github.com/justgook/goplatformer/pkg/resolv/v2"
)

type TagType = int64
type Level struct {
	Rooms    []*Room
	TileMaps []*image.Image
}

type Doors struct {
	N bool
	E bool
	S bool
	W bool
}

type Room struct {
	Layers    [][]int
	Doors     Doors
	W         int
	H         int
	Collision []*resolv.Object[TagType]
}

func (a *Level) Load(data []byte) error {
	var outputBuffer bytes.Buffer
	outputBuffer.Write(data)

	if err := gob.NewDecoder(&outputBuffer).Decode(a); err != nil {
		return err
	}

	return nil
}

func (a *Level) Save() ([]byte, error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)

	if err := encoder.Encode(*a); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
