package electrode

import (
	"testing"

	"paschen-pd/internal/paschen"
)

func TestParallelPlateBreakdown(t *testing.T) {
	pp := ParallelPlate{Pressure: 760, GapMm: 1}
	pd := pp.PD()
	if abs(pd-76) > 1e-9 {
		t.Errorf("PD = %g, want 76", pd)
	}
	if pp.BreakdownVoltage(paschen.DefaultAir()) <= 0 {
		t.Errorf("breakdown voltage should be positive")
	}
}

func TestSphereGapFactor(t *testing.T) {
	s := SphereGap{Pressure: 760, RadiusMm: 10, GapMm: 2}
	if s.UniformFieldFactor() <= 0 || s.UniformFieldFactor() > 1 {
		t.Errorf("field factor = %g, want in (0,1]", s.UniformFieldFactor())
	}
	// Large gap relative to radius -> strong non-uniformity.
	big := SphereGap{Pressure: 760, RadiusMm: 10, GapMm: 20}
	if big.UniformFieldFactor() != 0.6 {
		t.Errorf("field factor for gap>radius should be 0.6, got %g", big.UniformFieldFactor())
	}
}

func TestEvaluate(t *testing.T) {
	sc := []Scenario{
		{Name: "air-1mm", Pressure: 760, GapMm: 1, AppliedV: 1000, Params: paschen.DefaultAir()},
	}
	res := Evaluate(sc)
	if len(res) != 1 {
		t.Fatalf("expected 1 result")
	}
	if !res[0].Safe {
		t.Errorf("1 mm air at 1000 V should be safe, vb=%g", res[0].BreakdownV)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
