// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    animation, err := UnmarshalAnimation(bytes)
//    bytes, err = animation.Marshal()

package aseprite

import "encoding/json"

func UnmarshalAnimation(data []byte) (Animation, error) {
	var r Animation
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Animation) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Animation struct {
	Frames map[string]FrameValue `json:"frames"`
	Meta   Meta                  `json:"meta"`
}

type FrameValue struct {
	Frame            SpriteSourceSizeClass `json:"frame"`
	Rotated          bool                  `json:"rotated"`
	Trimmed          bool                  `json:"trimmed"`
	SpriteSourceSize SpriteSourceSizeClass `json:"spriteSourceSize"`
	SourceSize       Size                  `json:"sourceSize"`
	Duration         int64                 `json:"duration"`
}

type SpriteSourceSizeClass struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type Size struct {
	W int64 `json:"w"`
	H int64 `json:"h"`
}

type Meta struct {
	App       string        `json:"app"`
	Version   string        `json:"version"`
	Image     string        `json:"image"`
	Format    string        `json:"format"`
	Size      Size          `json:"size"`
	Scale     string        `json:"scale"`
	FrameTags []FrameTag    `json:"frameTags"`
	Layers    []Layer       `json:"layers"`
	Slices    []interface{} `json:"slices"`
}

type FrameTag struct {
	Name      string    `json:"name"`
	From      int       `json:"from"`
	To        int       `json:"to"`
	Direction Direction `json:"direction"`
	Color     Color     `json:"color"`
}

type Layer struct {
	Name      string `json:"name"`
	Opacity   int64  `json:"opacity"`
	BlendMode string `json:"blendMode"`
}

type Color string

const (
	The000000Ff Color = "#000000ff"
)

type Direction string

const (
	Forward Direction = "forward"
)

