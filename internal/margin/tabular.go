package margin

import (
	"fmt"
)

// Tabular renders several assessments into an aligned text table.
func Tabular(list []Assessment) string {
	out := "pd(Torr·cm)\tV_b(V)\tV_app(V)\tmargin%\tstate\n"
	for _, a := range list {
		state := "OK"
		if !a.Safe {
			state = "BREAK"
		}
		out += fmt.Sprintf("%.4g\t%.1f\t%.1f\t%.1f\t%s\n",
			a.PD, a.BreakdownV, a.AppliedV, a.MarginFraction*100, state)
	}
	return out
}

// Summarize returns a one-line summary of how many of n scenarios are safe.
func Summarize(list []Assessment) string {
	safe := SafeCount(list)
	return fmt.Sprintf("%d/%d 绝缘点安全", safe, len(list))
}

// BreakdownVoltageOf returns the breakdown voltage for a gas at a geometry using
// the margin package's gas lookup (air fallback).
func BreakdownVoltageOf(name string, pTorr, dMm float64) float64 {
	return AssessGas(name, pTorr, dMm, 0).BreakdownV
}
