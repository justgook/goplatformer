package room

import (
	"github.com/justgook/goplatformer/pkg/bin"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type objectType = resolv.Object[bin.TagType]
type FloatingPlatform struct {
	FloatingPlatform      *objectType
	FloatingPlatformTween *gween.Sequence
}

func NewFloatingPlatform(x, y, w, h float64) *FloatingPlatform {
	platformTag := bin.TagType(3)
	p := &FloatingPlatform{}
	var world struct {
		Space *resolv.Space[bin.TagType]

		FloatingPlatform      *objectType
		FloatingPlatformTween *gween.Sequence
	}
	// Create the floating platform.
	p.FloatingPlatform = resolv.NewObject(x, y, w, h, platformTag)

	// The floating platform moves using a *gween.Sequence sequence of tweens, moving it back and forth.
	p.FloatingPlatformTween = gween.NewSequence()
	world.FloatingPlatformTween.Add(
		gween.New(float32(p.FloatingPlatform.Y)-48, float32(p.FloatingPlatform.Y-128), 2, ease.Linear),
		gween.New(float32(p.FloatingPlatform.Y-128), float32(p.FloatingPlatform.Y)-48, 2, ease.Linear),
	)
	//world.Space.Add(world.FloatingPlatform)
	return p
}
func (p *FloatingPlatform) Update() {
	// Floating platform movement needs to be done before the player's movement update to make sure there's no space between its top and the player's bottom;
	// otherwise, an alternative might be to have the platform detect to see if the Player's resting on it, and if so, move the player up manually.
	y, _, seqDone := p.FloatingPlatformTween.Update(1.0 / 60.0)
	p.FloatingPlatform.Y = float64(y)
	if seqDone {
		p.FloatingPlatformTween.Reset()
	}
	p.FloatingPlatform.Update()
}
