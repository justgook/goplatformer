package goplatformer

import (
	_ "embed"
)

//go:embed asset/excel.ttf
var ExcelFont []byte

//go:embed asset/player/Combined.png
var AnimationPNG []byte

//go:embed asset/player/Combined.json
var AnimationJson []byte

//go:embed asset/tileset.png
var TilesetPng1 []byte

//go:embed asset/tileset2.png
var TilesetPng2 []byte

//go:embed asset/example.ldtk
var ExampleLevel []byte
