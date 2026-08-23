package sweep

import (
	"testing"

	"paschen-pd/internal/paschen"
)

func TestExtrema(t *testing.T) {
	c := Config{From: 0.1, To: 100, Points: 200, Params: paschen.DefaultAir()}
	e := Extrema(c)
	if e.Kind != "minimum" {
		t.Errorf("expected minimum, got %s", e.Kind)
	}
	if e.Voltage <= 0 {
		t.Errorf("extremum voltage = %g, want positive", e.Voltage)
	}
}

func TestFirstAbove(t *testing.T) {
	c := Config{From: 0.1, To: 100, Points: 200, Params: paschen.DefaultAir()}
	pd := FirstAbove(c, 1000)
	if pd <= 0 {
		t.Errorf("first-above-1000 pd = %g, want positive", pd)
	}
	if AboveThreshold(c, 1000) <= 0 {
		t.Errorf("expected some points above 1000 V")
	}
}
