package paschen

// Gas is a named gas with its Paschen coefficients.
type Gas struct {
	Name  string
	A     float64 // cm^-1 Torr^-1
	B     float64 // V cm^-1 Torr^-1
	Gamma float64 // secondary emission coefficient
}

// Params returns the gas coefficients as a Params value.
func (g Gas) Params() Params {
	return Params{A: g.A, B: g.B, Gamma: g.Gamma}
}

// KnownGases is a small registry of common gases and their empirical Paschen
// coefficients (A in cm^-1 Torr^-1, B in V cm^-1 Torr^-1, gamma dimensionless).
// Values are representative and intended for teaching/back-of-envelope use.
var KnownGases = map[string]Gas{
	"air":       {Name: "dry air", A: 15.0, B: 365.0, Gamma: 0.01},
	"argon":     {Name: "argon", A: 12.0, B: 175.0, Gamma: 0.01},
	"helium":    {Name: "helium", A: 3.0, B: 34.0, Gamma: 0.01},
	"hydrogen":  {Name: "hydrogen", A: 5.0, B: 130.0, Gamma: 0.01},
	"nitrogen":  {Name: "nitrogen", A: 12.0, B: 315.0, Gamma: 0.01},
	"co2":       {Name: "carbon dioxide", A: 20.0, B: 470.0, Gamma: 0.01},
}

// GasByName looks up a gas by key (matched case-insensitively against the
// registry keys). It returns the gas and true if found.
func GasByName(key string) (Gas, bool) {
	for k, g := range KnownGases {
		if equalFold(k, key) {
			return g, true
		}
	}
	return Gas{}, false
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

// BreakdownVoltageFor looks up a gas and returns its breakdown voltage at the
// given pressure-distance product (Torr·cm). Unknown gas falls back to air.
func BreakdownVoltageFor(name string, pd float64) float64 {
	g, ok := GasByName(name)
	if !ok {
		g = KnownGases["air"]
	}
	return g.Params().BreakdownVoltageValue(pd)
}

// MinimumFor returns the Paschen minimum (pd, V) for a named gas, defaulting to
// air when the name is unknown.
func MinimumFor(name string) (pdMin, vMin float64) {
	g, ok := GasByName(name)
	if !ok {
		g = KnownGases["air"]
	}
	return g.Params().Minimum()
}
