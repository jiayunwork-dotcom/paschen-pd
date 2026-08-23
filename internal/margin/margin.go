// Package margin analyses how close an applied voltage is to the breakdown
// threshold predicted by the Paschen model, given a gas and a gap geometry.
package margin

import (
	"paschen-pd/internal/paschen"
	"paschen-pd/internal/units"
)

// Assessment summarises the insulation margin for a single operating point.
type Assessment struct {
	PD             float64 // pressure-distance product, Torr·cm
	BreakdownV     float64 // predicted breakdown voltage, V
	AppliedV       float64 // applied voltage, V
	MarginFraction float64 // (Vb - Va)/Vb, negative means breakdown
	Safe           bool    // true when applied < breakdown
}

// Assess evaluates a single point: pressure (Torr), gap (mm), applied voltage.
func Assess(p params, pTorr, dMm, appliedV float64) Assessment {
	pd := units.PDProduct(pTorr, dMm)
	vb := p.BreakdownVoltageValue(pd)
	return Assessment{
		PD:             pd,
		BreakdownV:     vb,
		AppliedV:       appliedV,
		MarginFraction: p.MarginFraction(pd, appliedV),
		Safe:           vb > appliedV,
	}
}

// params aliases the paschen coefficient set so callers can pass it directly.
type params = paschen.Params

// AssessGas looks up a named gas and assesses the operating point. Unknown gas
// falls back to dry air.
func AssessGas(name string, pTorr, dMm, appliedV float64) Assessment {
	g, ok := paschen.GasByName(name)
	if !ok {
		g = paschen.KnownGases["air"]
	}
	return Assess(g.Params(), pTorr, dMm, appliedV)
}
