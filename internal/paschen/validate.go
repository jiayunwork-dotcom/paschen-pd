package paschen

import "errors"

// Validate reports whether the coefficient set is physically usable. A gas must
// have positive A, B and gamma; gamma must lie in (0, 1) for the ln(1+1/γ)
// term to be well defined and positive.
func (p Params) Validate() error {
	if p.A <= 0 {
		return errors.New("A must be positive")
	}
	if p.B <= 0 {
		return errors.New("B must be positive")
	}
	if p.Gamma <= 0 || p.Gamma >= 1 {
		return errors.New("gamma must be in (0, 1)")
	}
	return nil
}

// Sanitized returns a copy of p with non-positive coefficients replaced by the
// dry-air defaults, and validates the result.
func Sanitized(A, B, gamma float64) (Params, error) {
	p := Params{A: A, B: B, Gamma: gamma}
	if A <= 0 {
		p.A = 15.0
	}
	if B <= 0 {
		p.B = 365.0
	}
	if gamma <= 0 || gamma >= 1 {
		p.Gamma = 0.01
	}
	return p, p.Validate()
}

// IsBreakdownAboveThreshold reports whether V_b(pd) exceeds the given voltage.
func (p Params) IsBreakdownAboveThreshold(pd, threshold float64) bool {
	return p.BreakdownVoltageValue(pd) > threshold
}

// SafeGapDistance returns the largest mm gap that holds V at the given pressure.
func (p Params) SafeGapDistance(pTorr, v float64) float64 {
	g := p.SparkingDistance(pTorr, v)
	if g <= 0 {
		return 0
	}
	return g * 10.0 // cm -> mm
}
