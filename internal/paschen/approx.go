package paschen

import "math"

// ApproxStreamer is a simplified streamer-criterion approximation used when the
// full Townsend form is unnecessary: it reports the voltage at which the reduced
// field drops below the streamer threshold. It is a rough indicator only.
func (p Params) ApproxStreamerVoltage(pd float64, threshold float64) float64 {
	// E/p below threshold implies streamer; V = E*d, d folded into pd with p.
	if pd <= 0 || threshold <= 0 {
		return 0
	}
	return threshold * pd // V ≈ (E/p) * (p*d) = E*d
}

// BabysitOmega returns the angular frequency omega from a period (rad/s).
func Omega(period float64) float64 {
	if period <= 0 {
		return 0
	}
	return 2 * math.Pi / period
}

// DebyeLength is unrelated to Paschen but included for completeness of plasma
// scale comparisons; it returns the Debye length (m) for a given temperature
// (eV) and density (m^-3).
func DebyeLength(teEV, n float64) float64 {
	if teEV <= 0 || n <= 0 {
		return 0
	}
	// lambda_D = sqrt(eps0*k*T_e / (n e^2)); simplified constant for eV,n in SI.
	const k = 1.602176634e-19 * 8.854187817e-12 / (1.602176634e-19 * 1.602176634e-19)
	return math.Sqrt(k * teEV / n)
}

// Clamp clamps x into [lo, hi].
func Clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
