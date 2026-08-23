package paschen

import "fmt"

// GasComparisonRow is one row of a cross-gas comparison at a common pd.
type GasComparisonRow struct {
	Name    string
	PD      float64
	Voltage float64
}

// CompareGases evaluates several named gases at the same pressure-distance
// product and returns a row per gas, sorted by descending breakdown voltage.
// Unknown names fall back to dry air.
func CompareGases(pd float64, names ...string) []GasComparisonRow {
	rows := make([]GasComparisonRow, 0, len(names))
	for _, n := range names {
		g, ok := GasByName(n)
		if !ok {
			g = KnownGases["air"]
		}
		rows = append(rows, GasComparisonRow{
			Name:    g.Name,
			PD:      pd,
			Voltage: g.Params().BreakdownVoltageValue(pd),
		})
	}
	// Insertion sort by voltage descending.
	for i := 1; i < len(rows); i++ {
		for j := i; j > 0 && rows[j].Voltage > rows[j-1].Voltage; j-- {
			rows[j], rows[j-1] = rows[j-1], rows[j]
		}
	}
	return rows
}

// CompareReport renders the comparison as a short table string.
func CompareReport(pd float64, names ...string) string {
	rows := CompareGases(pd, names...)
	out := "gas\tV_b(V)\n"
	for _, r := range rows {
		out += fmt.Sprintf("%s\t%.2f\n", r.Name, r.Voltage)
	}
	return out
}
