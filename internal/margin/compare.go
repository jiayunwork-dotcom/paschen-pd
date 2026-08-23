package margin

import (
	"paschen-pd/internal/paschen"
	"paschen-pd/internal/units"
)

// MultiGas compares the insulation margin of several gases at one operating
// point and reports whether each holds. Unknown names fall back to air.
func MultiGas(pTorr, dMm, appliedV float64, names ...string) []Assessment {
	out := make([]Assessment, 0, len(names))
	for _, n := range names {
		out = append(out, AssessGas(n, pTorr, dMm, appliedV))
	}
	return out
}

// WorstOf returns the assessment with the smallest margin fraction among the
// list, or nil when the list is empty.
func WorstOf(list []Assessment) *Assessment {
	if len(list) == 0 {
		return nil
	}
	w := &list[0]
	for i := 1; i < len(list); i++ {
		if list[i].MarginFraction < w.MarginFraction {
			w = &list[i]
		}
	}
	return w
}

// SafeCount counts how many assessments are safe (applied < breakdown).
func SafeCount(list []Assessment) int {
	n := 0
	for _, a := range list {
		if a.Safe {
			n++
		}
	}
	return n
}

// pdForGas is a small helper returning the pd for a gas name at a given geometry.
func pdForGas(name string, pTorr, dMm float64) float64 {
	g, ok := paschen.GasByName(name)
	if !ok {
		g = paschen.KnownGases["air"]
	}
	_ = g
	return units.PDProduct(pTorr, dMm)
}
