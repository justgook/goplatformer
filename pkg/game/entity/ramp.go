package entity

import (
	"github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateRamp(ecs *ecs.ECS, obj *domain.Object) *donburi.Entry {
	ramp := Ramp.Spawn(ecs)
	components.Object.Set(ramp, obj)

	// We will construct the shape using a ConvexPolygon. It's essentially an elogated triangle, but with a "floor" afterwards,
	// ensuring the Player is always able to stand regardless of which ramp they're standing on.
	rampShape := resolv.NewConvexPolygon(
		0, 0,

		0, 0,
		2, 0, // The extra 2 pixels here make it so the Player doesn't get stuck for a frame or two when running up the ramp.
		obj.W-2, obj.H, // Same here; an extra 2 pixels makes it so that dismounting the ramp is nice and easy
		obj.W, obj.H,
		0, obj.H,
	)
	obj.SetShape(rampShape)

	return ramp
}
