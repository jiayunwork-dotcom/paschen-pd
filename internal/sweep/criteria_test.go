package sweep

import (
	"testing"

	"paschen-pd/internal/paschen"
)

func TestHoldable(t *testing.T) {
	c := Config{From: 1, To: 50, Points: 50, Params: paschen.DefaultAir()}
	// A tiny required voltage should be holdable across the window.
	if !Holdable(c, 10) {
		t.Errorf("window should hold 10 V")
	}
	if BandWidthAbove(c, 10) <= 0 {
		t.Errorf("band width above 10 V should be positive")
	}
	if ClassifyWindow(c, 10) != "safe-band" {
		t.Errorf("classify = %s, want safe-band", ClassifyWindow(c, 10))
	}
}

func TestCountAbove(t *testing.T) {
	c := Config{From: 1, To: 50, Points: 50, Params: paschen.DefaultAir()}
	if CountAbove(c, 1000) <= 0 {
		t.Errorf("expected some points above 1000 V")
	}
}
