package paschen

import (
	"math"
	"math/rand"
)

// UncertaintyConfig controls a Monte-Carlo propagation of coefficient errors.
type UncertaintyConfig struct {
	Samples      int     // number of samples
	RelA         float64 // relative uncertainty on A (e.g. 0.05 = 5%)
	RelB         float64 // relative uncertainty on B
	RelGamma     float64 // relative uncertainty on gamma
	Seed         int64
}

// UncertaintyBand estimates the spread of the breakdown voltage at a fixed pd
// by perturbing the coefficients within their relative tolerances. It returns
// the mean voltage and the (mean ± spread) bounds over the samples.
func (p Params) UncertaintyBand(pd float64, cfg UncertaintyConfig) (mean, low, high float64) {
	if cfg.Samples <= 1 {
		cfg.Samples = 64
	}
	rng := rand.New(rand.NewSource(cfg.Seed))
	var sum, mn, mx float64
	mn = math.MaxFloat64
	for i := 0; i < cfg.Samples; i++ {
		pa := p.A * (1 + cfg.RelA*(rng.Float64()*2-1))
		pb := p.B * (1 + cfg.RelB*(rng.Float64()*2-1))
		pg := p.Gamma * (1 + cfg.RelGamma*(rng.Float64()*2-1))
		if pg <= 0 {
			pg = p.Gamma
		}
		v := Params{A: pa, B: pb, Gamma: pg}.BreakdownVoltageValue(pd)
		sum += v
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	mean = sum / float64(cfg.Samples)
	return mean, mn, mx
}

// Sensitivity reports how strongly V depends on each coefficient at a given pd,
// using finite differences with the given step fraction.
type Sensitivity struct {
	DA float64 // dV/dA
	DB float64 // dV/dB
	DGamma float64 // dV/dGamma
}

// ComputeSensitivity returns the partial derivatives of the breakdown voltage
// with respect to A, B and gamma at the given pd.
func (p Params) ComputeSensitivity(pd, step float64) Sensitivity {
	if step <= 0 {
		step = 1e-3
	}
	base := p.BreakdownVoltageValue(pd)
	a2 := p
	a2.A += p.A * step
	b2 := p
	b2.B += p.B * step
	g2 := p
	g2.Gamma += p.Gamma * step
	vA := a2.BreakdownVoltageValue(pd)
	vB := b2.BreakdownVoltageValue(pd)
	vG := g2.BreakdownVoltageValue(pd)
	return Sensitivity{
		DA: (vA - base) / (p.A * step),
		DB: (vB - base) / (p.B * step),
		DGamma: (vG - base) / (p.Gamma * step),
	}
}
