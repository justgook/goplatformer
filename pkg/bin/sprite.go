package bin

import (
	"bytes"
	"encoding/gob"
)

type AnimatedSprite struct {
	Image []byte
	W     int
	H     int
	Data  AnimatedSpriteDataMap
}

type AnimatedSpriteDataMap = map[string][]AnimatedSpriteFrame
type AnimatedSpriteFrame struct {
	Duration int
	Layers   []AnimatedSpriteFrameLayer
}

type AnimatedSpriteFrameLayer struct {
	W  int
	H  int
	TX int
	TY int
	X0 int
	Y0 int
	X1 int
	Y1 int
}

func (a *AnimatedSprite) Load(data []byte) error {
	var outputBuffer bytes.Buffer
	outputBuffer.Write(data)

	if err := gob.NewDecoder(&outputBuffer).Decode(a); err != nil {
		return err
	}

	return nil
}

func (a *AnimatedSprite) Save() ([]byte, error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(*a); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
