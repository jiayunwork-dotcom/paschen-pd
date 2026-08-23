package margin

import (
	"bytes"
	"strings"
	"testing"
)

func TestExportCSV(t *testing.T) {
	list := []Assessment{
		AssessGas("air", 760, 1,   1000),
	}
	var buf bytes.Buffer
	if err := ExportCSV(&buf, list); err != nil {
		t.Fatalf("ExportCSV: %v", err)
	}
	if !strings.Contains(buf.String(), "pd_torr_cm") {
		t.Errorf("missing header")
	}
	a := list[0]
	if a.MarginalVoltage() <= 0 {
		t.Errorf("marginal voltage should be positive for safe case")
	}
	if a.Ratio() <= 1 {
		t.Errorf("ratio should exceed 1 for safe case")
	}
}
