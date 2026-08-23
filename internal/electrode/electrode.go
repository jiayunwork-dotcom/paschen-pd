// Package electrode models simple electrode geometries and maps them onto the
// Paschen breakdown model. The breakdown physics still lives in the paschen
// package; this package only translates geometry into the pd product and offers
// common electrode layouts.
package electrode

import (
	"paschen-pd/internal/paschen"
	"paschen-pd/internal/units"
)

// ParallelPlate describes two parallel planar electrodes separated by a gap.
type ParallelPlate struct {
	Pressure float64 // gas pressure, Torr
	GapMm    float64 // electrode separation, mm
}

// PD returns the pressure-distance product (Torr·cm) for the geometry.
func (p ParallelPlate) PD() float64 {
	return units.PDProduct(p.Pressure, p.GapMm)
}

// BreakdownVoltage returns the breakdown voltage using the given gas params.
func (p ParallelPlate) BreakdownVoltage(params paschen.Params) float64 {
	return params.BreakdownVoltageValue(p.PD())
}

// SphereGap approximates a sphere-to-plane gap using the same Paschen law (a
// rough but standard first estimate; real sphere gaps deviate near the tips).
type SphereGap struct {
	Pressure float64 // gas pressure, Torr
	RadiusMm float64 // sphere radius, mm
	GapMm    float64 // sphere-to-plane distance, mm
}

// PD returns the pressure-distance product (Torr·cm).
func (s SphereGap) PD() float64 {
	return units.PDProduct(s.Pressure, s.GapMm)
}

// BreakdownVoltage returns the estimated breakdown voltage.
func (s SphereGap) BreakdownVoltage(params paschen.Params) float64 {
	return params.BreakdownVoltageValue(s.PD())
}

// UniformFieldFactor returns the field-uniformity factor k for a sphere gap, a
// heuristic correction: k=1 when the gap is much smaller than the radius and
// drops as the gap approaches the radius.
func (s SphereGap) UniformFieldFactor() float64 {
	if s.RadiusMm <= 0 {
		return 1
	}
	ratio := s.GapMm / s.RadiusMm
	if ratio >= 1 {
		return 0.6 // strongly non-uniform
	}
	return 1 - 0.4*ratio
}

// CoaxialCylinder models an inner/outer cylindrical geometry (used for cables).
type CoaxialCylinder struct {
	Pressure   float64 // gas pressure, Torr
	InnerRadius float64 // inner conductor radius, mm
	OuterRadius float64 // outer conductor radius, mm
}

// EffectiveGap returns the characteristic gap (outer minus inner), mm.
func (c CoaxialCylinder) EffectiveGap() float64 {
	return c.OuterRadius - c.InnerRadius
}

// PD returns the pressure-distance product (Torr·cm) using the effective gap.
func (c CoaxialCylinder) PD() float64 {
	return units.PDProduct(c.Pressure, c.EffectiveGap())
}

// BreakdownVoltage returns the estimated breakdown voltage.
func (c CoaxialCylinder) BreakdownVoltage(params paschen.Params) float64 {
	return params.BreakdownVoltageValue(c.PD())
}
