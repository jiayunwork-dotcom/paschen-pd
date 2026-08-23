package paschen

import "testing"

func TestUncertaintyBandCentered(t *testing.T) {
	p := DefaultAir()
	mean, low, high := p.UncertaintyBand(76, UncertaintyConfig{
		Samples: 200, RelA: 0.05, RelB: 0.05, RelGamma: 0.05, Seed: 42,
	})
	if low <= 0 || high <= low {
		t.Errorf("band invalid: low=%g high=%g", low, high)
	}
	if mean <= 0 {
		t.Errorf("mean should be positive, got %g", mean)
	}
}

func TestSensitivityMagnitude(t *testing.T) {
	p := DefaultAir()
	s := p.ComputeSensitivity(76, 1e-3)
	if s.DB <= 0 {
		t.Errorf("dV/dB should be positive, got %g", s.DB)
	}
	if abs(s.DA) < 0 {
		t.Errorf("dV/dA should be defined, got %g", s.DA)
	}
}
