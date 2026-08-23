package api

import "paschen-pd/internal/paschen"

// BreakdownRequest is the JSON body for /api/breakdown.
//
// The distance d and pressure p define the gap; pd is computed in Torr·cm
// (d is supplied in mm, p in Torr, so pd = p * d / 10).
type BreakdownRequest struct {
	P     float64 `json:"p"`     // pressure, Torr
	D     float64 `json:"d"`     // gap distance, mm
	A     float64 `json:"A"`     // Townsend A coefficient (0 => default air)
	B     float64 `json:"B"`     // Townsend B coefficient (0 => default air)
	Gamma float64 `json:"gamma"` // secondary emission coefficient (0 => default)
}

// BreakdownResponse carries the computed breakdown voltage and supporting data.
type BreakdownResponse struct {
	PD            float64 `json:"pd_torr_cm"` // pressure-distance product, Torr·cm
	Voltage       float64 `json:"voltage"`    // breakdown voltage, V
	PDMin         float64 `json:"pd_min"`     // pd at the Paschen minimum
	VMin          float64 `json:"v_min"`      // minimum breakdown voltage, V
	A             float64 `json:"A"`
	B             float64 `json:"B"`
	Gamma         float64 `json:"gamma"`
	Breakdownable bool    `json:"breakdownable"`
	LeftBranch    bool    `json:"left_branch"` // true when pd < pd_min (rising branch)
}

// SweepRequest is the JSON body for /api/sweep.
type SweepRequest struct {
	P      float64 `json:"p"`      // fixed pressure, Torr (for constant-pressure scan)
	DMin   float64 `json:"d_min"`  // smallest gap, mm
	DMax   float64 `json:"d_max"`  // largest gap, mm
	Points int     `json:"points"` // number of samples
	Log    bool    `json:"log"`    // logarithmic spacing
	A      float64 `json:"A"`
	B      float64 `json:"B"`
	Gamma  float64 `json:"gamma"`
}

// SweepResponse carries the sampled curve (each point already in Torr·cm).
type SweepResponse struct {
	Points []SweepPoint `json:"points"`
	A      float64      `json:"A"`
	B      float64      `json:"B"`
	Gamma  float64      `json:"gamma"`
}

// SweepPoint is one sample of the curve.
type SweepPoint struct {
	PD      float64 `json:"pd_torr_cm"`
	Voltage float64 `json:"voltage"`
}

// defaultIfZero returns the gas coefficients, substituting air defaults when a
// supplied value is non-positive.
func defaultIfZero(req BreakdownRequest) (a, b, g float64) {
	a, b, g = req.A, req.B, req.Gamma
	if a <= 0 {
		a = 15.0
	}
	if b <= 0 {
		b = 365.0
	}
	if g <= 0 {
		g = 0.01
	}
	return a, b, g
}

// paramsFrom builds a paschen.Params from a request, applying defaults.
func paramsFrom(req BreakdownRequest) paschen.Params {
	a, b, g := defaultIfZero(req)
	return paschen.Params{A: a, B: b, Gamma: g}
}
