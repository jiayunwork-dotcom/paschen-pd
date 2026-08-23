// Package units converts between the pressure and length units used in gas
// discharge physics. The Paschen model works in Torr for pressure and cm for
// length, but inputs often arrive in Pascal, bar, or mm. See pressure.go and
// length.go for the individual conversions.
package units

import "math"

// PDProduct computes the pressure-distance product in Torr·cm from a pressure
// in Torr and a gap distance in mm.
func PDProduct(pTorr, dMm float64) float64 {
	return pTorr * LengthFromMm(dMm)
}

// RoundTo rounds x to n decimal places.
func RoundTo(x float64, n int) float64 {
	scale := math.Pow(10, float64(n))
	return math.Round(x*scale) / scale
}
