package paschen

import "testing"

func TestBreakdownVoltageKnownAir(t *testing.T) {
	p := DefaultAir()
	// At 760 Torr and 1 mm gap -> pd = 76 Torr·cm.
	// V = 365*76 / (ln(15*76) - ln(ln(101))) = 27740 / (7.037 - 1.529) = 5035 V
	got := p.BreakdownVoltage(76.0)
	if got < 4800 || got > 5200 {
		t.Errorf("air 760Torr/1mm breakdown = %g V, want ~5035", got)
	}
}

func TestBreakdownZeroBelowRegime(t *testing.T) {
	p := DefaultAir()
	// Very small pd -> denominator becomes non-positive -> no breakdown.
	if v := p.BreakdownVoltage(1e-4); v != 0 {
		t.Errorf("tiny pd should give V=0, got %g", v)
	}
}

func TestMinimumIsConsistent(t *testing.T) {
	p := DefaultAir()
	pdMin, vMin := p.Minimum()
	// Optimal when ln(A·pd*) - ln(ln(1+1/γ)) = 1, with C=ln(ln(101))=1.529,
	// so ln(A·pd*) = C+1 -> pd* = exp(2.529)/15 ≈ 0.836 Torr·cm.
	if pdMin < 0.7 || pdMin > 1.0 {
		t.Errorf("pd_min = %g, want ~0.84 Torr·cm", pdMin)
	}
	// The minimum must be lower than the voltage at a much larger pd.
	if p.BreakdownVoltage(100) <= vMin {
		t.Errorf("minimum voltage %g should be below far-field %g", vMin, p.BreakdownVoltage(100))
	}
	// Re-evaluating the minimum point reproduces vMin.
	if got := p.BreakdownVoltage(pdMin); abs(got-vMin) > 1e-6 {
		t.Errorf("V(pd_min)=%g != vMin=%g", got, vMin)
	}
}

func TestReducedFieldPositive(t *testing.T) {
	p := DefaultAir()
	if rf := p.ReducedField(76.0); rf <= 0 {
		t.Errorf("reduced field must be positive, got %g", rf)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
