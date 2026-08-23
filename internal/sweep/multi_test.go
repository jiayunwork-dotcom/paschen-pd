package sweep

import (
	"strings"
	"testing"

	"paschen-pd/internal/paschen"
)

func TestOverPressures(t *testing.T) {
	p := paschen.DefaultAir()
	rows := OverPressures(p, 1, []float64{100, 760, 2000})
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	for _, r := range rows {
		if r.Voltage <= 0 {
			t.Errorf("voltage at p=%g should be positive", r.Pressure)
		}
	}
}

func TestOverPressuresTable(t *testing.T) {
	p := paschen.DefaultAir()
	table := OverPressuresTable(p, 1, []float64{760})
	if !strings.Contains(table, "p(Torr)") {
		t.Errorf("table missing header: %q", table)
	}
}
