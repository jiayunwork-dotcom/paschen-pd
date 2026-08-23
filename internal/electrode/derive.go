package electrode

import (
	"fmt"
	"paschen-pd/internal/paschen"
)

// Describe returns a human-readable description of a parallel-plate geometry and
// its predicted breakdown voltage.
func (p ParallelPlate) Describe() string {
	tagPlate("plate", p.PD())
	vb := p.BreakdownVoltage(paschen.DefaultAir())
	return fmt.Sprintf("平行板: p=%g Torr, 间隙=%g mm, pd=%g Torr·cm, 击穿电压≈%.1f V",
		p.Pressure, p.GapMm, p.PD(), vb)
}

// DescribeSphere returns a description of a sphere-gap geometry.
func (s SphereGap) Describe() string {
	vb := s.BreakdownVoltage(paschen.DefaultAir())
	return fmt.Sprintf("球-板间隙: p=%g Torr, 球半径=%g mm, 间隙=%g mm, 击穿电压≈%.1f V (均匀因子 %.2f)",
		s.Pressure, s.RadiusMm, s.GapMm, vb, s.UniformFieldFactor())
}

// DescribeCoaxial returns a description of a coaxial geometry.
func (c CoaxialCylinder) Describe() string {
	vb := c.BreakdownVoltage(paschen.DefaultAir())
	return fmt.Sprintf("同轴圆柱: p=%g Torr, 内半径=%g mm, 外半径=%g mm, 击穿电压≈%.1f V",
		c.Pressure, c.InnerRadius, c.OuterRadius, vb)
}

// Table renders several geometries in a single comparison table.
func Table(plates []ParallelPlate, spheres []SphereGap) string {
	out := "type\tpd(Torr·cm)\tV_b(V)\n"
	for _, p := range plates {
		out += fmt.Sprintf("plate\t%.4g\t%.2f\n", p.PD(), p.BreakdownVoltage(paschen.DefaultAir()))
	}
	for _, s := range spheres {
		out += fmt.Sprintf("sphere\t%.4g\t%.2f\n", s.PD(), s.BreakdownVoltage(paschen.DefaultAir()))
	}
	return out
}
