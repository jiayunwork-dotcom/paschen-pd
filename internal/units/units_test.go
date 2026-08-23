package units

import "testing"

func TestPressureConversions(t *testing.T) {
	// 760 Torr == 1 atm.
	if got := PressureFromAtm(1); abs(got-760) > 1e-6 {
		t.Errorf("1 atm = %g Torr, want 760", got)
	}
	// 1 bar ≈ 750.06 Torr.
	if got := PressureFromBar(1); got < 749 || got > 752 {
		t.Errorf("1 bar = %g Torr, want ~750", got)
	}
	if got := PressureToPa(760); abs(got-101325) > 50 {
		t.Errorf("760 Torr = %g Pa, want ~101325", got)
	}
}

func TestLengthConversions(t *testing.T) {
	if got := LengthFromMm(10); got != 1 {
		t.Errorf("10 mm = %g cm, want 1", got)
	}
	if got := LengthFromM(1); got != 100 {
		t.Errorf("1 m = %g cm, want 100", got)
	}
}

func TestPDProduct(t *testing.T) {
	// 760 Torr * 1 mm -> pd = 76 Torr·cm.
	if got := PDProduct(760, 1); abs(got-76) > 1e-9 {
		t.Errorf("PDProduct(760,1) = %g, want 76", got)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
