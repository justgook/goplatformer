package particle_test

import (
	"github.com/justgook/goplatformer/pkg/particle"
	"image/color"
	"testing"
	"time"
)

func TestParticle_System(t *testing.T) {
	s := particle.NewSystem()
	p := particle.NewParticle(s)
	assertEqual(t, p.System(), s)
}

func TestParticle_Update(t *testing.T) {
	sys := particle.NewSystem()

	sys.MaxParticles = 1

	sys.LifetimeOverTime = func(d time.Duration, delta time.Duration) time.Duration {
		return 1500 * time.Millisecond
	}

	sys.DataOverLifetime = func(old any, t particle.NormalizedDuration, delta time.Duration) any {
		return "data"
	}

	sys.EmissionPositionOverTime = func(d time.Duration, delta time.Duration) particle.Vector {
		return particle.Vector{X: 17, Y: 23}
	}

	sys.VelocityOverLifetime = func(p *particle.Particle, t particle.NormalizedDuration, delta time.Duration) particle.Vector {
		return particle.Vector{X: 3, Y: 5}
	}

	sys.ScaleOverLifetime = func(p *particle.Particle, t particle.NormalizedDuration, delta time.Duration) particle.Vector {
		return particle.Vector{X: 7, Y: 11}
	}

	sys.ColorOverLifetime = func(p *particle.Particle, t particle.NormalizedDuration, delta time.Duration) color.Color {
		return color.RGBA{R: 0x12, G: 0x23, B: 0x34, A: 0x45}
	}

	sys.RotationOverLifetime = func(p *particle.Particle, t particle.NormalizedDuration, delta time.Duration) float64 {
		return 0.123
	}

	updateCalled := false
	sys.UpdateFunc = func(part *particle.Particle, t particle.NormalizedDuration, delta time.Duration) {
		updateCalled = true
	}

	deathCalled := false
	sys.DeathFunc = func(p *particle.Particle) {
		deathCalled = true
	}

	sys.Spawn(1)

	now := time.Now()
	sys.Update(now)

	var part *particle.Particle

	sys.ForEachParticle(func(p *particle.Particle, t particle.NormalizedDuration, delta time.Duration) {
		part = p
	}, now)

	assertEqual(t, part.Data(), "data")
	assertEqual(t, part.Position(), particle.Vector{X: 17, Y: 23})
	assertEqual(t, part.Velocity(), particle.Vector{X: 3, Y: 5})
	assertEqual(t, part.Scale(), particle.Vector{X: 7, Y: 11})
	assertColor(t, part.Color(), color.RGBA{R: 0x12, G: 0x23, B: 0x34, A: 0x45})
	assertEqual(t, part.Angle(), 0.0)
	assertEqual(t, part.Lifetime(), 1500*time.Millisecond)
	assertEqual(t, updateCalled, true)

	now = now.Add(1 * time.Second)
	sys.Update(now)

	assertEqual(t, part.Position(), particle.Vector{X: 17, Y: 23}.Add(particle.Vector{X: 3, Y: 5}))
	assertEqual(t, part.Angle(), 0.123)

	now = now.Add(1 * time.Second)
	sys.Update(now)
	assertEqual(t, deathCalled, true)
}

func TestParticle_Kill(t *testing.T) {

	sys := particle.NewSystem()

	sys.MaxParticles = 1

	sys.LifetimeOverTime = func(d time.Duration, delta time.Duration) time.Duration {
		return 10 * time.Second
	}

	sys.Spawn(1)

	now := time.Now()
	sys.Update(now)

	var part *particle.Particle

	sys.ForEachParticle(func(p *particle.Particle, t particle.NormalizedDuration, delta time.Duration) {
		part = p
	}, now)

	part.Kill()

	now = now.Add(1 * time.Second)
	sys.Update(now)

	assertEqual(t, sys.NumParticles(), 0)
}
