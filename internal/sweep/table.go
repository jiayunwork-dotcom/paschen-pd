package sweep

import (
	"fmt"
	"strings"
)

// Table renders a human-readable breakdown-voltage table for the sampled curve.
// It prints pd (Torr·cm) and the corresponding breakdown voltage in volts.
func Table(c Config) string {
	pts := Run(c)
	var b strings.Builder
	b.WriteString("pd(Torr·cm)\tV(V)\n")
	for _, p := range pts {
		b.WriteString(fmt.Sprintf("%.6g\t%.3f\n", p.PD, p.Voltage))
	}
	return b.String()
}

// Stats summarises the sampled curve: the minimum voltage found within the
// scanned window and the pd at which it occurs.
func Stats(c Config) (pdAtMin, vMin float64) {
	pts := Run(c)
	vMin = pts[0].Voltage
	pdAtMin = pts[0].PD
	for _, p := range pts {
		if p.Voltage < vMin {
			vMin = p.Voltage
			pdAtMin = p.PD
		}
	}
	return pdAtMin, vMin
}

// PeakBetween finds the largest gap (in terms of pd) that still breaks down
// below a ceiling voltage, within the scanned window. It returns the bounding
// pd and whether such a point was found.
func PeakBetween(c Config, ceiling float64) (pd float64, ok bool) {
	pts := Run(c)
	for i := len(pts) - 1; i >= 0; i-- {
		if pts[i].Voltage <= ceiling {
			return pts[i].PD, true
		}
	}
	return -1, false
}

// Envelope returns the (pd, V) pair with the largest voltage in the window.
func Envelope(c Config) (pd, vMax float64) {
	pts := Run(c)
	vMax = pts[0].Voltage
	pd = pts[0].PD
	for _, p := range pts {
		if p.Voltage > vMax {
			vMax = p.Voltage
			pd = p.PD
		}
	}
	return pd, vMax
}

