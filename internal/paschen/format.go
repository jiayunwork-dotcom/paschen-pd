package paschen

import "fmt"

// FormatBreakdown returns a compact one-line description of the breakdown at a
// given pd.
func (p Params) FormatBreakdown(pd float64) string {
	v := p.BreakdownVoltageValue(pd)
	return fmt.Sprintf("pd=%.4g Torr·cm -> V_b=%.2f V", pd, v)
}

// FormatMinimum returns a description of the Paschen minimum.
func (p Params) FormatMinimum() string {
	pd, v := p.Minimum()
	return fmt.Sprintf("Paschen最小: pd=%.4g Torr·cm, V_min=%.2f V", pd, v)
}

// FormatComparisonRow renders a (name, voltage) pair for cross-gas output.
func FormatComparisonRow(name string, v float64) string {
	return fmt.Sprintf("%-12s %.2f V", name, v)
}
