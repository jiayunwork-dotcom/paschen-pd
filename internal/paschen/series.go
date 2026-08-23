package paschen

import "math"

// SeriesPoint is one point of a voltage/distance sweep at fixed pressure.
type SeriesPoint struct {
	D    float64 // gap distance, cm
	PD   float64 // pressure-distance product, Torr·cm
	V    float64 // breakdown voltage, V
}

// Series evaluates the breakdown voltage across a set of gap distances (cm) at a
// fixed pressure (Torr). It is the basis for plotting V(d) at constant p.
func (p Params) Series(pTorr float64, gapsCm []float64) []SeriesPoint {
	out := make([]SeriesPoint, 0, len(gapsCm))
	for _, d := range gapsCm {
		pd := pTorr * d
		out = append(out, SeriesPoint{
			D:    d,
			PD:   pd,
			V:    p.BreakdownVoltageValue(pd),
		})
	}
	return out
}

// MaxHoldableVoltage returns the largest voltage that the geometry can sustain:
// it is the breakdown voltage at the maximum gap in the supplied set (since for
// pd beyond the minimum the curve is monotonic increasing).
func (p Params) MaxHoldableVoltage(pTorr float64, gapsCm []float64) float64 {
	best := 0.0
	for _, d := range gapsCm {
		v := p.BreakdownVoltageValue(pTorr * d)
		if v > best {
			best = v
		}
	}
	return best
}

// MinVoltageGap returns the gap (cm) within the set that yields the smallest
// breakdown voltage (closest to the Paschen minimum).
func (p Params) MinVoltageGap(pTorr float64, gapsCm []float64) (gap, vMin float64) {
	bestGap := gapsCm[0]
	bestV := p.BreakdownVoltageValue(pTorr * bestGap)
	for _, d := range gapsCm[1:] {
		v := p.BreakdownVoltageValue(pTorr * d)
		if v < bestV {
			bestV = v
			bestGap = d
		}
	}
	return bestGap, bestV
}

// DerivativedVdpd returns the numerical derivative dV/d(pd) at a given pd using a
// centred difference with the supplied step.
func (p Params) DerivativedVdpd(pd, step float64) float64 {
	if step <= 0 {
		step = 1e-4
	}
	vp := p.BreakdownVoltageValue(pd + step)
	vm := p.BreakdownVoltageValue(pd - step)
	return (vp - vm) / (2 * step)
}

// IsNearMinimum reports whether pd is within tol of the Paschen minimum.
func (p Params) IsNearMinimum(pd, tol float64) bool {
	pdMin, _ := p.Minimum()
	return math.Abs(pd-pdMin) <= tol
}
