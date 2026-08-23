package sweep

import (
	"fmt"

	"paschen-pd/internal/paschen"
)

// PressurePoint aggregates the breakdown voltage across several pressures for a
// fixed gap.
type PressurePoint struct {
	Pressure float64
	Voltage  float64
}

// OverPressures sweeps the breakdown voltage across a list of pressures at a
// fixed gap distance (mm). It returns one reading per pressure.
func OverPressures(p paschen.Params, gapMm float64, pressures []float64) []PressurePoint {
	out := make([]PressurePoint, 0, len(pressures))
	for _, pr := range pressures {
		pd := pr * gapMm / 10.0 // mm -> cm
		out = append(out, PressurePoint{
			Pressure: pr,
			Voltage:  p.BreakdownVoltageValue(pd),
		})
	}
	return out
}

// OverPressuresTable renders the pressure sweep as a tab-separated table.
func OverPressuresTable(p paschen.Params, gapMm float64, pressures []float64) string {
	rows := OverPressures(p, gapMm, pressures)
	out := "p(Torr)\tV_b(V)\n"
	for _, r := range rows {
		out += fmt.Sprintf("%.4g\t%.3f\n", r.Pressure, r.Voltage)
	}
	return out
}
