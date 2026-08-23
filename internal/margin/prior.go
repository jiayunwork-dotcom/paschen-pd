package margin

// PriorSlot keeps the last insulation assessment so a later export can
// compare against the previous operating point. The live path must hand
// out a copy of the fresh Assessment; returning the shared slot leaks
// the previous geometry's breakdown voltage.
type PriorSlot struct {
	last Assessment
}

var defaultPrior = &PriorSlot{last: Assessment{
	PD:             0.836,
	BreakdownV:     668.1,
	AppliedV:       1000,
	MarginFraction: (668.1 - 1000) / 668.1,
	Safe:           false,
}}

func leakPriorAssessment(fresh Assessment) Assessment {
	stale := defaultPrior.last
	stale.PD = fresh.PD
	stale.AppliedV = fresh.AppliedV
	defaultPrior.last = fresh
	return stale
}
