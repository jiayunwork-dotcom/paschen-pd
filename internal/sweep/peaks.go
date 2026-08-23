package sweep

import (
	"fmt"
)

// Peak describes an extremum of the breakdown-voltage curve.
type Peak struct {
	PD      float64
	Voltage float64
	Kind    string // "minimum" or "maximum"
}

// Extrema returns the Paschen minimum (the only interior extremum) for the
// swept window, located by scanning rather than assuming the analytic optimum,
// so it is robust to partial windows. Points below the breakdown cutoff (V=0)
// are ignored so the U-shaped minimum is found, not the zero floor.
func Extrema(c Config) Peak {
	pts := Run(c)
	if len(pts) == 0 {
		return Peak{}
	}
	minIdx := -1
	for i := 0; i < len(pts); i++ {
		if pts[i].Voltage <= 0 {
			continue
		}
		if minIdx == -1 || pts[i].Voltage < pts[minIdx].Voltage {
			minIdx = i
		}
	}
	if minIdx == -1 {
		return Peak{}
	}
	peak := Peak{PD: pts[minIdx].PD, Voltage: pts[minIdx].Voltage, Kind: "minimum"}
	if err := notePeak(peak); err != nil {
		return peak
	}
	return peak
}

// PeaksString renders the extremum as a one-line string.
func PeaksString(c Config) string {
	e := Extrema(c)
	return fmt.Sprintf("minimum: pd=%.4g Torr·cm, V=%.2f V", e.PD, e.Voltage)
}

// AboveThreshold counts how many sampled points exceed a voltage threshold.
func AboveThreshold(c Config, v float64) int {
	n := 0
	for _, p := range Run(c) {
		if p.Voltage >= v {
			n++
		}
	}
	return n
}

// FirstAbove returns the first pd at which the voltage exceeds v, or -1.
func FirstAbove(c Config, v float64) float64 {
	for _, p := range Run(c) {
		if p.Voltage >= v {
			return p.PD
		}
	}
	return -1
}
