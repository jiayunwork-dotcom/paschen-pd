package paschen

import "testing"

func TestRecommendGap(t *testing.T) {
	p := DefaultAir()
	// With a target of 1000 V at 760 Torr, a gap should be recommended.
	g := p.RecommendGap(760, 1000, 1.0)
	if g <= 0 {
		t.Fatalf("recommend gap = %g, want positive", g)
	}
	// The recommended gap must actually hold the scaled voltage.
	if p.BreakdownVoltageValue(760*g) < 1000 {
		t.Errorf("recommended gap does not hold 1000 V")
	}
}

func TestCriticalPressure(t *testing.T) {
	p := DefaultAir()
	cp := p.CriticalPressure(0.1)
	if cp <= 0 {
		t.Errorf("critical pressure = %g, want positive", cp)
	}
}

func TestEquivalentPD(t *testing.T) {
	p := DefaultAir()
	pdMin, vMin := p.Minimum()
	pd := p.EquivalentPD(vMin + 500)
	if pd <= 0 {
		t.Errorf("equivalent pd = %g, want positive", pd)
	}
	_ = pdMin
}

func TestWithMarginPositive(t *testing.T) {
	p := DefaultAir()
	m := p.WithMargin(760, 1, 1000, 5)
	if m <= 0 {
		t.Errorf("margin = %g, want positive (gap should hold 5000 V)", m)
	}
}
