package particle_test

import (
	"github.com/justgook/goplatformer/pkg/particle"
	"math"
	"testing"
)

func TestVector_Magnitude(t *testing.T) {
	assertEqual(t,
		particle.Vector{X: 17, Y: 23}.Magnitude(),
		math.Sqrt(17*17+23*23),
	)
}

func TestVector_TryNormalize(t *testing.T) {

	v := particle.Vector{X: 17, Y: 23}
	m := v.Magnitude()

	norm, ok := v.TryNormalize()
	assertEqual(t, norm.X, v.X/m)
	assertEqual(t, norm.Y, v.Y/m)
	assertEqual(t, norm.Magnitude(), 1.0)
	assertEqual(t, ok, true)

	v = particle.Vector{}
	norm, ok = v.TryNormalize()
	assertEqual(t, v, norm)
	assertEqual(t, ok, false)

}

func TestVector_Add(t *testing.T) {
	v1 := particle.Vector{X: 17, Y: 23}
	v2 := particle.Vector{X: 5, Y: 7}
	assertEqual(t, v1.Add(v2), particle.Vector{X: v1.X + v2.X, Y: v1.Y + v2.Y})
}

func TestVector_Multiply(t *testing.T) {
	assertEqual(t, particle.Vector{X: 17, Y: 23}.Multiply(3), particle.Vector{X: 17 * 3, Y: 23 * 3})
}
