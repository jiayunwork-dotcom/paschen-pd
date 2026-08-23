package paschen

import "testing"

func TestSummarize(t *testing.T) {
	p := DefaultAir()
	s := p.Summarize(76)
	if s.Voltage <= 0 {
		t.Errorf("summary voltage = %g, want positive", s.Voltage)
	}
	if s.ReducedField <= 0 {
		t.Errorf("summary reduced field should be positive")
	}
	if s.PDMin <= 0 {
		t.Errorf("summary pdMin should be positive")
	}
	if len(s.Describe()) == 0 {
		t.Errorf("describe should not be empty")
	}
}

func TestBulkSummary(t *testing.T) {
	p := DefaultAir()
	out := p.BulkSummary([]float64{1, 10, 100})
	if len(out) != 3 {
		t.Fatalf("expected 3 summaries, got %d", len(out))
	}
}
