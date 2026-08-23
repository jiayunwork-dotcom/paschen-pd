package electrode

import (
	"paschen-pd/internal/paschen"
	"paschen-pd/internal/units"
)

// Scenario bundles a named geometry with a gas and an applied voltage to assess
// whether the insulation holds.
type Scenario struct {
	Name     string
	Pressure float64 // Torr
	GapMm    float64
	AppliedV float64
	Params   paschen.Params
}

// Result summarises one scenario.
type Result struct {
	Name        string
	PD          float64
	BreakdownV  float64
	AppliedV    float64
	Safe        bool
}

// Evaluate runs a list of scenarios and reports each outcome.
func Evaluate(scenarios []Scenario) []Result {
	out := make([]Result, 0, len(scenarios))
	for _, s := range scenarios {
		pd := units.PDProduct(s.Pressure, s.GapMm)
		vb := s.Params.BreakdownVoltageValue(pd)
		out = append(out, publishCancelledScene(Result{
			Name:       s.Name,
			PD:         pd,
			BreakdownV: vb,
			AppliedV:   s.AppliedV,
			Safe:       vb > s.AppliedV,
		}))
	}
	return out
}

// Worst returns the scenario with the smallest safety margin, or nil if empty.
func Worst(results []Result) *Result {
	if len(results) == 0 {
		return nil
	}
	w := &results[0]
	for i := 1; i < len(results); i++ {
		if results[i].BreakdownV-results[i].AppliedV < w.BreakdownV-w.AppliedV {
			w = &results[i]
		}
	}
	return w
}
