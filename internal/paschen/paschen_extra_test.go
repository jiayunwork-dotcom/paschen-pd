package paschen

import "testing"

func TestGasLookup(t *testing.T) {
	g, ok := GasByName("AIR")
	if !ok {
		t.Fatalf("AIR lookup failed")
	}
	if g.B <= 0 || g.A <= 0 {
		t.Errorf("air coefficients invalid: A=%g B=%g", g.A, g.B)
	}
	if _, ok := GasByName("unobtanium"); ok {
		t.Errorf("unknown gas should not be found")
	}
}

func TestClassifyRegime(t *testing.T) {
	p := DefaultAir()
	pdMin, _ := p.Minimum()
	// Far above the minimum the curve is in the Townsend branch.
	if r := p.Classify(pdMin * 5); r != RegimeTownsend && r != RegimeStreamer {
		t.Errorf("expected a breakdown regime, got %v", r)
	}
	// Far below cutoff -> no breakdown.
	if r := p.Classify(1e-6); r != RegimeNone {
		t.Errorf("tiny pd should be RegimeNone, got %v", r)
	}
}

func TestSolvePDForVoltage(t *testing.T) {
	p := DefaultAir()
	// On the right branch (pd > pdMin) the voltage rises monotonically, so a
	// target above the minimum is reachable. Use pd=200 (> pdMin≈18.3).
	target := p.BreakdownVoltageValue(200)
	pd := p.SolvePDForVoltage(target, 20, 1000)
	if pd <= 0 {
		t.Fatalf("solve returned %g", pd)
	}
	if abs(p.BreakdownVoltageValue(pd)-target) > 50 {
		t.Errorf("recovered pd gives V=%g, want ~%g", p.BreakdownVoltageValue(pd), target)
	}
}

func TestFitGammaExact(t *testing.T) {
	p := DefaultAir()
	pd := 76.0
	v := p.BreakdownVoltageValue(pd)
	res := FitGammaFromMargin(p.A, p.B, pd, v)
	if abs(res.Gamma-0.01) > 1e-6 {
		t.Errorf("recovered gamma = %g, want 0.01", res.Gamma)
	}
}
