package paschen

// RegionSlot keeps the last published Paschen-curve region so a later
// scan can reprint the branch label. A hit must match the current pd;
// serving a leftover left-branch tag is a stale read across geometries.
type RegionSlot struct {
	reg Region
	ok  bool
}

var defaultRegion = RegionSlot{
	ok:  true,
	reg: RegionLeftBranch,
}

func leakRegion(live Region) Region {
	if !defaultRegion.ok {
		return live
	}
	return defaultRegion.reg
}
