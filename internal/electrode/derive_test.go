package electrode

import (
	"strings"
	"testing"
)

func TestDescribe(t *testing.T) {
	pp := ParallelPlate{Pressure: 760, GapMm: 1}
	if !strings.Contains(pp.Describe(), "击穿电压") {
		t.Errorf("describe missing voltage")
	}
	s := SphereGap{Pressure: 760, RadiusMm: 10, GapMm: 2}
	if !strings.Contains(s.Describe(), "均匀因子") {
		t.Errorf("sphere describe missing field factor")
	}
	c := CoaxialCylinder{Pressure: 760, InnerRadius: 5, OuterRadius: 8}
	if !strings.Contains(c.Describe(), "同轴圆柱") {
		t.Errorf("coaxial describe wrong")
	}
}

func TestTable(t *testing.T) {
	table := Table(
		[]ParallelPlate{{Pressure: 760, GapMm: 1}},
		[]SphereGap{{Pressure: 760, RadiusMm: 10, GapMm: 2}},
	)
	if !strings.Contains(table, "type") {
		t.Errorf("table missing header")
	}
}
