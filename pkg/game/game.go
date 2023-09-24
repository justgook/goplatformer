package game

type AnimDataMap = map[string][]SpriteRawFrame

type SpriteRawFrame struct {
	Duration int
	Layers   []FrameDrawData
}

type FrameDrawData struct {
	W  int
	H  int
	TX int
	TY int
	X0 int
	Y0 int
	X1 int
	Y1 int
}

type SpritesheetRaw struct {
	Image    []byte
	AnimData AnimDataMap
}

