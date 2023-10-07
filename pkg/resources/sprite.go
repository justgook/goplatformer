package resources

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

