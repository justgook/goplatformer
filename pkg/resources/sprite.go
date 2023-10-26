package resources

import (
	"image"
	"time"
)

type AnimatedSprite struct {
	Image []byte
	W     int
	H     int
	Tags  []string
	Anim  map[string]AnimatedSpriteAnimData

	Data AnimatedSpriteDataMap
}
type AnimatedSpriteAnimData struct {
	Frames    []*image.Rectangle
	Durations []time.Duration
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

func (a *AnimatedSprite) SpriteSheet() error {

	return nil
}
