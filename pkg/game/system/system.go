package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game/state"
	"github.com/justgook/goplatformer/pkg/util"
)

var _ state.Scene = (*Systems)(nil)

type Systems []state.Scene

func (s Systems) Init() {
	for i := range s {
		s[i].Init()
	}
}

func (s Systems) Draw(screen *ebiten.Image) {
	for i := range s {
		s[i].Draw(screen)
	}
}

func (s Systems) Update(aa *state.GameState) error {
	var err error
	for i := range s {
		if err = s[i].Update(aa); err != nil {
			return util.Catch(err)
		}
	}
	return nil
}

func (s Systems) Terminate() {
	for i := range s {
		s[i].Terminate()
	}
}
