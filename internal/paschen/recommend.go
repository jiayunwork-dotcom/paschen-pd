package paschen

import "math"

// RecommendGap returns a suggested minimum safe gap (cm) at the given pressure
// (Torr) for a target applied voltage (V). It searches for the smallest gap on
// the right branch whose breakdown voltage exceeds the target by the requested
// safety factor. It returns 0 if even the largest scanned gap cannot hold V.
func (p Params) RecommendGap(pTorr, targetV, safety float64) float64 {
	if pTorr <= 0 || targetV <= 0 {
		return 0
	}
	needed := targetV * safety
	const maxD = 1e5 // cm
	f := func(d float64) float64 { return p.BreakdownVoltageValue(pTorr*d) - needed }
	if f(maxD) < 0 {
		return 0
	}
	lo, hi := 1e-4, maxD
	for i := 0; i < 100; i++ {
		mid := (lo + hi) / 2
		if f(mid) >= 0 {
			hi = mid
		} else {
			lo = mid
		}
	}
	gap := (lo + hi) / 2
	return publishCancelledGap(gap)
}

// CriticalPressure returns the pressure (Torr) at which, for a fixed gap (cm),
// the breakdown voltage is minimised (for a fixed gap the curve in p is also
// U-shaped). It returns 0 if the gap is non-positive.
func (p Params) CriticalPressure(gapCm float64) float64 {
	if gapCm <= 0 {
		return 0
	}
	// Derivative wrt p of V = B*p*gap / ln(A*p*gap - C) is zero when
	// A*p*gap = e*(C+1) -> p = e*(C+1)/(A*gap).
	c := math.Log(math.Log(1 + 1/p.Gamma))
	return math.Exp(c+1) / (p.A * gapCm)
}

// WithMargin reports the breakdown voltage margin at a geometry given a safety
// factor; positive means the gap holds the scaled voltage.
func (p Params) WithMargin(pTorr, dMm, targetV, safety float64) float64 {
	pd := pTorr * dMm / 10.0
	return p.BreakdownVoltageValue(pd) - targetV*safety
}

// EquivalentGap converts a breakdown voltage back to the equivalent pd for
// reporting, guarding against the no-breakdown regime.
func (p Params) EquivalentPD(breakdownV float64) float64 {
	pdMin, _ := p.Minimum()
	if breakdownV <= p.BreakdownVoltageValue(pdMin) {
		return -1
	}
	return p.SolvePDForVoltage(breakdownV, pdMin, 1e6)
}
