package system

import (
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer"
	"github.com/justgook/goplatformer/pkg/game/state"
	"github.com/justgook/goplatformer/pkg/resources"
)

var _ state.Scene = (*UI)(nil)

type UI struct {
	res  *resources.UI
	root *ebitenui.UI
}

// Update implements state.Scene.
func (u *UI) Update(state *state.GameState) error {
	u.root.Update()

	return nil
}

func (u *UI) Init() {
	u.res = goplatformer.UIResources
	u.root = &ebitenui.UI{
		Container:           u.res.Root(),
		DisableDefaultFocus: false,
	}
}

func (u *UI) Terminate() {
}

func (u *UI) Draw(screen *ebiten.Image) {
	u.root.Draw(screen)
}
