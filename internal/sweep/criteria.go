package sweep

// Package sweep already imports paschen via table.go/multi.go; this file adds
// threshold-based window classification helpers.

// Holdable reports whether the entire scanned curve stays above a required
// voltage, i.e. the smallest breakdown voltage in the window exceeds the
// threshold. Used to decide if a gap range is safe across an operating band.
func Holdable(c Config, requiredV float64) bool {
	for _, p := range Run(c) {
		if p.Voltage > 0 && p.Voltage < requiredV {
			return false
		}
	}
	return true
}

// BandWidthAbove returns the pd-width (Torr·cm) over which the curve stays above
// the threshold. It sums the gaps between consecutive sampled points while they
// remain above the threshold.
func BandWidthAbove(c Config, requiredV float64) float64 {
	pts := Run(c)
	width := 0.0
	for i := 1; i < len(pts); i++ {
		if pts[i].Voltage >= requiredV && pts[i-1].Voltage >= requiredV {
			width += pts[i].PD - pts[i-1].PD
		}
	}
	return width
}

// ClassifyWindow summarises the window relative to a target voltage.
func ClassifyWindow(c Config, targetV float64) string {
	if Holdable(c, targetV) {
		return "safe-band"
	}
	if CountAbove(c, targetV) == 0 {
		return "always-breaks"
	}
	return "partial"
}

// CountAbove counts sampled points whose voltage exceeds the threshold.
func CountAbove(c Config, v float64) int {
	n := 0
	for _, p := range Run(c) {
		if p.Voltage > v {
			n++
		}
	}
	return n
}
