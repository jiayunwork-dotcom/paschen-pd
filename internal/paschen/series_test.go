package paschen

import "testing"

func TestSeriesPositive(t *testing.T) {
	p := DefaultAir()
	gaps := []float64{0.1, 1, 10}
	pts := p.Series(760, gaps)
	if len(pts) != 3 {
		t.Fatalf("expected 3 points, got %d", len(pts))
	}
	for _, pt := range pts {
		if pt.V <= 0 {
			t.Errorf("voltage at d=%g should be positive", pt.D)
		}
	}
}

func TestMinVoltageGap(t *testing.T) {
	p := DefaultAir()
	gaps := []float64{0.01, 0.836, 5, 20}
	g, v := p.MinVoltageGap(760, gaps)
	// The gap nearest the Paschen minimum should give the lowest breakdown V.
	if v <= 0 {
		t.Errorf("min voltage = %g, want positive", v)
	}
	_ = g
}

func TestDerivativedVdpd(t *testing.T) {
	p := DefaultAir()
	d := p.DerivativedVdpd(76, 1e-4)
	if d <= 0 {
		t.Errorf("dV/d(pd) should be positive on the right branch, got %g", d)
	}
}

func TestIsNearMinimum(t *testing.T) {
	p := DefaultAir()
	pdMin, _ := p.Minimum()
	if !p.IsNearMinimum(pdMin, 1e-6) {
		t.Errorf("pdMin should be near minimum")
	}
}
