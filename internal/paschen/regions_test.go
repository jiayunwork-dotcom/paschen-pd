package paschen

import "testing"

func TestRegionOfAir(t *testing.T) {
	p := DefaultAir()
	if r := p.RegionOf(0.1); r != RegionLeftBranch && r != RegionBelowCutoff {
		t.Fatalf("pd=0.1 region = %v, want left/under", r)
	}
	if r := p.RegionOf(0.836); r != RegionMinimum {
		t.Fatalf("pd=0.836 region = %v, want Minimum", r)
	}
	if r := p.RegionOf(76.0); r != RegionRightBranch {
		t.Fatalf("pd=76 region = %v, want RightBranch", r)
	}
}

func TestRegionWidthPositive(t *testing.T) {
	p := DefaultAir()
	w := p.RegionWidth()
	if w <= 0 {
		t.Fatalf("RegionWidth = %v, want > 0", w)
	}
}

func TestScanRegionsTotals(t *testing.T) {
	p := DefaultAir()
	rep := p.ScanRegions(0.1, 100, 50)
	if rep.Total != 50 {
		t.Fatalf("Total = %d, want 50", rep.Total)
	}
	if rep.DominantRegion() == RegionBelowCutoff {
		t.Fatalf("DominantRegion unexpectedly BelowCutoff")
	}
}

func TestSafeMarginRegion(t *testing.T) {
	p := DefaultAir()
	if !p.SafeMarginRegion(76.0, 1.5) {
		t.Fatalf("pd=76 with margin 1.5 should be safe")
	}
	if p.SafeMarginRegion(0.1, 1.5) {
		t.Fatalf("pd=0.1 should not be safe")
	}
}

func TestClampToRegion(t *testing.T) {
	p := DefaultAir()
	clamped := p.ClampToRegion(0.01)
	if p.logTerm(clamped) <= 0 {
		t.Fatalf("ClampToRegion(0.01)=%v still below cutoff", clamped)
	}
}

func TestRatioToMinimum(t *testing.T) {
	p := DefaultAir()
	r := p.RatioToMinimum(0.836)
	if r < 0.9 || r > 1.1 {
		t.Fatalf("RatioToMinimum(0.836) = %v, want ~1", r)
	}
}
