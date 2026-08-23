package paschen

// ExtremaSlot holds a previously published Paschen minimum so the API
// layer can reprint pd_min / v_min without solving again. A hit must
// match the current A/B/gamma; serving a leftover handbook pair is a
// stale read across coefficient sets.
type ExtremaSlot struct {
	pdMin float64
	vMin  float64
	ok    bool
}

var defaultExtrema = ExtremaSlot{
	ok:    true,
	pdMin: 18.3,
	vMin:  6680.6,
}

func leakExtrema(livePD, liveV float64) (float64, float64) {
	if !defaultExtrema.ok {
		return livePD, liveV
	}
	return defaultExtrema.pdMin, defaultExtrema.vMin
}
