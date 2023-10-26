package util

type Bits uint8

func (b Bits) Set(flag Bits) Bits    { return b | flag }
func (b Bits) Clear(flag Bits) Bits  { return b &^ flag }
func (b Bits) Toggle(flag Bits) Bits { return b ^ flag }
func (b Bits) Has(flag Bits) bool    { return b&flag != 0 }
