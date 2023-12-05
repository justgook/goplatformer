package particle_test

import (
	"github.com/justgook/goplatformer/pkg/particle"
	"testing"
	"time"
)

func TestParticleSystem_Reset(t *testing.T) {

	sys := particle.NewSystem()

	sys.MaxParticles = 1

	sys.LifetimeOverTime = func(d time.Duration, delta time.Duration) time.Duration {
		return 10 * time.Second
	}

	sys.Spawn(1)

	now := time.Now()
	sys.Update(now)

	sys.Reset()

	got := sys.NumParticles()
	want := 0
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestParticleSystem_Update_SpawnMoreAfterKill(t *testing.T) {
	sys := particle.NewSystem()

	sys.MaxParticles = 1

	sys.EmissionRateOverTime = func(d time.Duration, delta time.Duration) float64 {
		return 1.0
	}

	sys.LifetimeOverTime = func(d time.Duration, delta time.Duration) time.Duration {
		return 10 * time.Second
	}

	sys.Spawn(1)

	now := time.Now()
	sys.Update(now)

	killCalled := false
	sys.UpdateFunc = func(p *particle.Particle, t particle.NormalizedDuration, delta time.Duration) {
		if t > 0 {
			killCalled = true

			p.Kill()
		}
	}

	now = now.Add(1 * time.Second)
	sys.Update(now)

	t.Run("particle count", func(t *testing.T) {
		assertEqual(t, sys.NumParticles(), 1)
	})
	t.Run("killCalled", func(t *testing.T) {
		assertEqual(t, killCalled, true)
	})

}

func TestParticleSystem_Spawn(t *testing.T) {

	sys := particle.NewSystem()

	sys.MaxParticles = 1

	sys.Spawn(1)

	now := time.Now()
	sys.Update(now)

	t.Run("particle count", func(t *testing.T) {
		assertEqual(t, sys.NumParticles(), 1)
	})
}

func TestNormalizedDuration_Duration(t *testing.T) {
	t.Run("duration", func(t *testing.T) {
		assertEqual(t, particle.NormalizedDuration(0.2).Duration(5000*time.Millisecond), 1000*time.Millisecond)
	})
}
