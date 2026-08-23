package paschen

import "math"

// CriticalField returns the reduced field E/p (V cm^-1 Torr^-1) implied by the
// breakdown voltage V at the given pd, i.e. (V/d) normalised by pressure with
// d folded into pd. At breakdown this equals the Townsend condition.
func (p Params) CriticalField(pd, v float64) float64 {
	if pd <= 0 {
		return 0
	}
	return v / pd
}

// SparkingDistance returns the largest gap (in cm) that can still be held at the
// applied voltage V without breakdown, for the given pressure (Torr). It solves
// V = B·p·d / ln(A·p·d - C) for d by scanning, returning the root on the
// right branch (d > d_min). It returns 0 if V exceeds the Paschen minimum's
// voltage (impossible to hold) or if no root exists in [eps, 1e6] cm.
func (p Params) SparkingDistance(pTorr, v float64) float64 {
	if pTorr <= 0 || v <= 0 {
		return 0
	}
	// Convert to pd-space where pd = pTorr * d(cm). Search d in [eps, big].
	const eps = 1e-4
	const big = 1e5
	// f(d) = V_breakdown(pTorr*d) - v; find where it equals 0 on rising branch.
	f := func(d float64) float64 {
		return p.BreakdownVoltageValue(pTorr*d) - v
	}
	lo, hi := eps, big
	if f(lo) > 0 || f(hi) < 0 {
		return 0
	}
	for i := 0; i < 100; i++ {
		mid := (lo + hi) / 2
		if math.Abs(f(mid)) < 1e-6 {
			return mid
		}
		if f(lo)*f(mid) < 0 {
			hi = mid
		} else {
			lo = mid
		}
	}
	return (lo + hi) / 2
}

// PaschenNumber returns the scaled coordinate pd / pd_min used to locate a point
// on the U-shaped curve (1 at the minimum).
func PaschenNumberFor(p Params, pd float64) float64 {
	pdMin, _ := p.Minimum()
	if pdMin <= 0 {
		return 0
	}
	return pd / pdMin
}

// Region reports whether the point (pd, V) lies above, on, or below the
// breakdown curve. "above" means the gap holds (safe); "below" means breakdown.
func (p Params) Region(pd, v float64) string {
	vb := p.BreakdownVoltageValue(pd)
	switch {
	case vb <= 0:
		return "no-breakdown-regime"
	case v < vb:
		return "holds"
	case v > vb:
		return "breakdown"
	default:
		return "on-curve"
	}
}
