// Package paschen implements the Paschen law for gas breakdown voltage.
//
// The breakdown voltage of a gas in a parallel-plate gap follows the Paschen
// curve: there is a minimum voltage below which no breakdown can occur,
// regardless of how small the gap gets. The model used here is the classic
// Townsend form
//
//	V_b(pd) = B * p * d / ( ln(A * p * d) - ln( ln(1 + 1/gamma) ) )
//
// with p·d expressed in Torr·cm. The constants A, B and the secondary
// emission coefficient gamma depend on the gas; for dry air the commonly
// quoted values are A = 15 (cm^-1 Torr^-1), B = 365 (V cm^-1 Torr^-1) and
// gamma = 0.01.
package paschen

import "math"

// Params holds the gas coefficients used by the Paschen model.
type Params struct {
	A     float64 // first Townsend coefficient prefactor, cm^-1 Torr^-1
	B     float64 // second Townsend coefficient prefactor, V cm^-1 Torr^-1
	Gamma float64 // secondary electron emission coefficient
}

// DefaultAir returns the coefficient set for dry air.
func DefaultAir() Params {
	return Params{A: 15.0, B: 365.0, Gamma: 0.01}
}

// logTerm returns ln(A * pd) - ln(ln(1 + 1/gamma)), the denominator of the
// Paschen equation. It is positive only for pd above the cutoff where the
// gas can sustain breakdown.
func (p Params) logTerm(pd float64) float64 {
	return math.Log(p.A*pd) - math.Log(math.Log(1+1/p.Gamma))
}

// BreakdownVoltage returns the breakdown voltage (V) for a product of pressure
// p (Torr) and gap distance d (cm). The argument pd is already in Torr·cm.
// It returns 0 when pd is below the regime where the gas can break down.
func (p Params) BreakdownVoltage(pd float64) float64 {
	return p.BreakdownVoltageValue(pd)
}

// BreakdownVoltageValue is the same computation as BreakdownVoltage; it is kept
// as a plain-float method so helpers can call it without allocating.
func (p Params) BreakdownVoltageValue(pd float64) float64 {
	lt := p.logTerm(pd)
	if lt <= 0 {
		return 0
	}
	return p.B * pd / lt
}

// ReducedField returns the reduced electric field E/p (V cm^-1 Torr^-1) at
// breakdown, equal to V_b / d with d expressed in cm.
func (p Params) ReducedField(pd float64) float64 {
	d := pd // pd already carries the d component in cm units
	if d <= 0 {
		return 0
	}
	return p.BreakdownVoltage(pd) / d
}

// Minimum returns the pressure-distance product (Torr·cm) at which the
// breakdown voltage is smallest, together with that minimum voltage. Minimising
// V(pd) = B·pd/(ln(A·pd) - C) with C = ln(ln(1+1/γ) gives the optimum when the
// denominator equals 1, i.e. ln(A·pd*) = C + 1.
func (p Params) Minimum() (pdMin float64, vMin float64) {
	c := math.Log(math.Log(1 + 1/p.Gamma))
	if c <= 0 {
		return 0, 0
	}
	pdMin = math.Exp(c+1) / p.A
	vMin = p.BreakdownVoltageValue(pdMin)
	return leakExtrema(pdMin, vMin)
}
