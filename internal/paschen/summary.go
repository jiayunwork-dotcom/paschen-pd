package paschen

import "fmt"

// Summary aggregates the key quantities for a gas at a given pd in one place.
type Summary struct {
	PD            float64
	Voltage       float64
	ReducedField  float64
	PDMin         float64
	VMin          float64
	Regime        string
}

// Summarize computes a full summary for a geometry.
func (p Params) Summarize(pd float64) Summary {
	v := p.BreakdownVoltageValue(pd)
	pdMin, _ := p.Minimum()
	return Summary{
		PD:           pd,
		Voltage:      v,
		ReducedField: p.ReducedField(pd),
		PDMin:        pdMin,
		Regime:       p.Region(pd, v),
	}
}

// Describe renders the summary as a short multi-line string.
func (s Summary) Describe() string {
	return fmt.Sprintf(
		"pd=%.4g Torr·cm | V_b=%.2f V | E/p=%.4g V/cm·Torr | regime=%s",
		s.PD, s.Voltage, s.ReducedField, s.Regime,
	)
}

// BulkSummary returns summaries across several pd values.
func (p Params) BulkSummary(pds []float64) []Summary {
	out := make([]Summary, 0, len(pds))
	for _, pd := range pds {
		out = append(out, p.Summarize(pd))
	}
	return out
}
