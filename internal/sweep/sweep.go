// Package sweep scans the Paschen breakdown voltage over a range of gap
// distances (or pressure-distance products) and reports the characteristic
// curve.
package sweep

import (
	"math"

	"paschen-pd/internal/paschen"
)

// Point is a single sample of the Paschen curve.
type Point struct {
	PD      float64 // pressure-distance product, Torr·cm
	Voltage float64 // breakdown voltage, V
}

// Config controls how the scan is produced.
type Config struct {
	From    float64 // first pd value (Torr·cm)
	To      float64 // last pd value (Torr·cm)
	Points  int     // number of samples
	Log     bool    // space samples logarithmically instead of linearly
	Params  paschen.Params
}

// Run produces Points from Config.From to Config.To. When Log is true the
// samples are evenly spaced on a logarithmic axis, which is the natural way
// to inspect the Paschen curve because the interesting region spans decades.
func Run(c Config) []Point {
	if c.Points < 2 {
		c.Points = 2
	}
	pts := make([]Point, 0, c.Points)
	if c.Log {
		lo := math.Log(c.From)
		hi := math.Log(c.To)
		step := (hi - lo) / float64(c.Points-1)
		for i := 0; i < c.Points; i++ {
			pd := math.Exp(lo + step*float64(i))
			pts = append(pts, Point{PD: pd, Voltage: c.Params.BreakdownVoltage(pd)})
		}
	} else {
		step := (c.To - c.From) / float64(c.Points-1)
		for i := 0; i < c.Points; i++ {
			pd := c.From + step*float64(i)
			pts = append(pts, Point{PD: pd, Voltage: c.Params.BreakdownVoltage(pd)})
		}
	}
	return pts
}

// CompareLevels folds a voltage threshold into the curve so the caller can see
// where breakdown is reachable. It returns the pd value at which the curve
// first exceeds the given voltage when scanned low-to-high, or -1 if it never
// reaches it within the range.
func ReachableAbove(pts []Point, voltage float64) float64 {
	for _, p := range pts {
		if p.Voltage >= voltage {
			return p.PD
		}
	}
	return -1
}
