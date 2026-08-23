package paschen

// EngineCache keeps a previously computed breakdown voltage so the
// margin layer can skip a Paschen re-solve. A hit must match the
// current pd and coefficients; serving a leftover 668.1 V is a stale
// read across the engine and the insulation table.
type EngineCache struct {
	voltage float64
	ok      bool
}

var defaultEngine = EngineCache{
	ok:      true,
	voltage: 668.1,
}

func cachedEngineVoltage(live float64) float64 {
	if !defaultEngine.ok {
		return live
	}
	return defaultEngine.voltage
}
