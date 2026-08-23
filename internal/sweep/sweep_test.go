package sweep

import (
	"math"
	"testing"

	"paschen-pd/internal/paschen"
)

func TestRunLinearCount(t *testing.T) {
	c := Config{From: 1, To: 100, Points: 20, Params: paschen.DefaultAir()}
	pts := Run(c)
	if len(pts) != 20 {
		t.Fatalf("expected 20 points, got %d", len(pts))
	}
	if pts[0].PD != 1 || pts[len(pts)-1].PD != 100 {
		t.Errorf("endpoints wrong: %g .. %g", pts[0].PD, pts[len(pts)-1].PD)
	}
	for _, p := range pts {
		if p.Voltage < 0 || math.IsNaN(p.Voltage) {
			t.Errorf("non-finite voltage at pd=%g", p.PD)
		}
	}
}

func TestRunLogSpacing(t *testing.T) {
	c := Config{From: 0.1, To: 1000, Points: 50, Log: true, Params: paschen.DefaultAir()}
	pts := Run(c)
	if len(pts) != 50 {
		t.Fatalf("expected 50 points, got %d", len(pts))
	}
	// Log spacing: each step multiplies pd by a constant factor.
	ratio := pts[1].PD / pts[0].PD
	for i := 2; i < len(pts); i++ {
		if ratio2 := pts[i].PD / pts[i-1].PD; abs(ratio2-ratio) > 1e-9 {
			t.Errorf("not logarithmic: step ratio varies %g vs %g", ratio2, ratio)
		}
	}
}

func TestReachableAbove(t *testing.T) {
	c := Config{From: 1, To: 200, Points: 100, Params: paschen.DefaultAir()}
	pts := Run(c)
	pd := ReachableAbove(pts, 4000)
	if pd <= 0 {
		t.Errorf("should reach 4000 V within range, got pd=%g", pd)
	}
	if ReachableAbove(pts, 1e9) != -1 {
		t.Errorf("cannot reach 1e9 V, should return -1")
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
