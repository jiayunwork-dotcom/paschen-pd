package paschen

import "math"

// ElectricFieldAt returns the breakdown field E = V_b / d (V/cm) at the given
// pressure-distance product (Torr·cm), with d already folded into pd.
func (p Params) ElectricFieldAt(pd float64) float64 {
	if pd <= 0 {
		return 0
	}
	return p.BreakdownVoltageValue(pd) / pd
}

// Alpha computes the first Townsend coefficient alpha (cm^-1) from the reduced
// field E/p (V cm^-1 Torr^-1). It is the microscopic relation
//
//	alpha = A * p * exp(-B / (E/p))
//
// that underlies the Paschen law. Here p is the pressure in Torr; the reduced
// field already absorbs the p factor, so callers pass the combined quantity.
func Alpha(a, b, reducedField float64) float64 {
	if reducedField <= 0 {
		return 0
	}
	return a * math.Exp(-b/reducedField)
}

// TownsendAlphaAt returns the first Townsend coefficient (cm^-1) at the
// breakdown condition for the given pd. At breakdown alpha·d = ln(1 + 1/gamma),
// which this function reproduces.
func (p Params) TownsendAlphaAt(pd float64) float64 {
	ef := p.ElectricFieldAt(pd)
	if ef <= 0 {
		return 0
	}
	return Alpha(p.A, p.B, ef*pd/math.Max(pd, 1e-9))
}
