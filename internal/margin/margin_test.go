package margin

import "testing"

func TestAssessSafe(t *testing.T) {
	// Dry air at 760 Torr, 1 mm gap => pd=76, V_b≈11445 V. Applied 1000 V is safe.
	a := AssessGas("air", 760, 1, 1000)
	if !a.Safe {
		t.Errorf("expected safe at 1000 V, got margin %g", a.MarginFraction)
	}
	if abs(a.BreakdownV-5035) > 300 {
		t.Errorf("breakdown voltage = %g, want ~5035", a.BreakdownV)
	}
}

func TestAssessUnsafe(t *testing.T) {
	// Tiny gap under vacuum-like pd -> no breakdown, applied voltage exceeds.
	a := AssessGas("air", 760, 1, 100000)
	if a.Safe {
		t.Errorf("100 kV should exceed breakdown, got margin %g", a.MarginFraction)
	}
}

func TestReportContainsState(t *testing.T) {
	a := AssessGas("air", 760, 1, 1000)
	if len(a.Report()) == 0 {
		t.Errorf("report should not be empty")
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
