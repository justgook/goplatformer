package sprite

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/bin"
)

type Animated struct {
	Raw *bin.AnimatedSprite
}

func (a *Animated) Draw(target *ebiten.Image) {

}

func (a *Animated) SetName(name string) {

}
func (a *Animated) Update() {

}

func (a *Animated) FromRaw(input *bin.AnimatedSprite) {
	a.Raw = input
}

func (a *Animated) Load(data []byte) error {
	delme := &bin.AnimatedSprite{}
	if err := delme.Load(data); err != nil {
		return err
	}
	a.FromRaw(delme)

	return nil
}
