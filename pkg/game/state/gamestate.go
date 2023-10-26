package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/pcg"
)

type DeviceInfo struct {
	ScreenWidth  int
	ScreenHeight int
}

type GameState struct {
	Rand         *pcg.PCG32
	DeviceInfo   *DeviceInfo
	SceneManager *SceneManager
	Input        *InputManager
}

func (g *GameState) Init(w, h int, seed uint64) {
	g.Rand = pcg.NewPCG32().Seed(seed, 1)
	g.SceneManager = &SceneManager{
		current:         nil,
		next:            nil,
		transitionCount: 0,
		transitionFrom:  ebiten.NewImage(w, h),
		transitionTo:    ebiten.NewImage(w, h),
	}

	g.Input = &InputManager{}
	g.Input.Init()

}
func (g *GameState) SetScene(scene Scene) {
	s := g.SceneManager
	if s.current == nil {
		s.current = scene
		s.current.Init(g)

		s.runner.Update = s.current.Update
		s.runner.Draw = s.current.Draw
	} else {
		s.next = scene
		s.next.Init(g)
		s.transitionCount = transitionMaxCount

		s.runner.Update = transitionUpdate
		s.runner.Draw = g.transitionDraw
	}
}

func transitionUpdate(g *GameState) error {
	s := g.SceneManager
	s.transitionCount--
	if s.transitionCount > 0 {
		return nil
	}
	s.current.Terminate()
	s.current = s.next
	s.next = nil
	s.runner.Update = s.current.Update
	s.runner.Draw = s.current.Draw

	return nil
}

func (g *GameState) Update() error {
	g.Input.Update()
	return g.SceneManager.runner.Update(g)
}

func (g *GameState) transitionDraw(r *ebiten.Image) {
	s := g.SceneManager

	s.transitionFrom.Clear()
	s.current.Draw(s.transitionFrom)

	s.transitionTo.Clear()
	s.next.Draw(s.transitionTo)

	r.DrawImage(s.transitionFrom, nil)

	alpha := 1 - float32(s.transitionCount)/float32(transitionMaxCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(alpha)
	r.DrawImage(s.transitionTo, op)
}

func (g *GameState) Draw(r *ebiten.Image) {
	g.SceneManager.runner.Draw(r)
}
