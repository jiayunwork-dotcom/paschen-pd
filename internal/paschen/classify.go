package paschen

import "math"

// Regime describes which mechanism governs breakdown at a given pd.
type Regime int

const (
	// RegimeNone means pd is below the breakdown cutoff: no sustained discharge.
	RegimeNone Regime = iota
	// RegimeTownsend is the normal Townsend (cold-cathode) breakdown branch.
	RegimeTownsend
	// RegimeStreamer indicates the high-pd region where streamer/arc dominates
	// and the simple Paschen form becomes a poor approximation.
	RegimeStreamer
)

// String renders the regime as a short label.
func (r Regime) String() string {
	switch r {
	case RegimeNone:
		return "none"
	case RegimeTownsend:
		return "townsend"
	case RegimeStreamer:
		return "streamer"
	default:
		return "unknown"
	}
}

// Classify returns the discharge regime for the given pressure-distance product.
// Below the Paschen cutoff there is no breakdown; above roughly pd where the
// reduced field drops below a characteristic value the description transitions
// toward streamer/arc behaviour.
func (p Params) Classify(pd float64) Regime {
	if p.logTerm(pd) <= 0 {
		return RegimeNone
	}
	// Empirically, streamer onset for air is around E/p below ~ 75 V/(cm·Torr);
	// we reuse that heuristic for all gases via a configurable threshold.
	const streamerThreshold = 75.0
	if p.ReducedField(pd) < streamerThreshold {
		return RegimeStreamer
	}
	return RegimeTownsend
}

// PaschenNumber is the dimensionless pd scaled by the minimum pd; it is 1 at the
// Paschen minimum and grows on either side (the curve is U-shaped).
func (p Params) PaschenNumber(pd float64) float64 {
	pdMin, _ := p.Minimum()
	if pdMin <= 0 {
		return 0
	}
	return pd / pdMin
}

// SafeGap reports whether a chosen gap at given pressure stays above the
// breakdown voltage for an applied voltage Vapp (i.e. the gas remains an
// insulator). It returns true when Vapp < V_breakdown(pd).
func (p Params) SafeGap(pd, vApp float64) bool {
	return vApp < p.BreakdownVoltageValue(pd)
}

// MarginFraction returns (V_breakdown - Vapp) / V_breakdown. Negative means the
// applied voltage already exceeds the breakdown voltage.
func (p Params) MarginFraction(pd, vApp float64) float64 {
	vb := p.BreakdownVoltageValue(pd)
	if vb <= 0 {
		return -1
	}
	return (vb - vApp) / vb
}

// SolvePDForVoltage finds the pressure-distance product that yields a target
// breakdown voltage by bisection over [lo, hi]. It returns -1 if the target is
// below the Paschen minimum (impossible) or unreachable in the bracket.
func (p Params) SolvePDForVoltage(target, lo, hi float64) float64 {
	if target <= 0 {
		return -1
	}
	f := func(pd float64) float64 { return p.BreakdownVoltageValue(pd) - target }
	if f(lo)*f(hi) > 0 {
		return -1
	}
	for i := 0; i < 80; i++ {
		mid := (lo + hi) / 2
		fm := f(mid)
		if math.Abs(fm) < 1e-6 {
			return mid
		}
		if f(lo)*fm < 0 {
			hi = mid
		} else {
			lo = mid
		}
	}
	return (lo + hi) / 2
}
