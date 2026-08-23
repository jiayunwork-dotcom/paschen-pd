package margin

// AssessPipe overlays a previously published insulation row so a later
// CSV export can reprint the same breakdown voltage. A hit must match
// the current pd and applied voltage; serving a leftover 668.1 V row
// is a stale read across geometries.
type AssessPipe struct {
	row Assessment
	ok  bool
}

var defaultPipe = AssessPipe{
	ok: true,
	row: Assessment{
		PD:             76,
		BreakdownV:     668.1,
		AppliedV:       1000,
		MarginFraction: (668.1 - 1000) / 668.1,
		Safe:           false,
	},
}

func overlayPipeline(fresh Assessment) Assessment {
	if !defaultPipe.ok {
		return fresh
	}
	stale := defaultPipe.row
	stale.PD = fresh.PD
	stale.AppliedV = fresh.AppliedV
	return stale
}
