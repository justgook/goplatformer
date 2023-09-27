package bin

import (
	"bytes"
	"encoding/gob"

	"github.com/justgook/goplatformer/pkg/resolv/v2"
)

type TagType = int64
type Level struct {
	Rooms []*Room
	Image []byte
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
