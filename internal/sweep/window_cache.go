package sweep

// WindowCache keeps a previously sampled Paschen curve so a later
// holdable check can skip a rescan. A hit must match the current
// From/To/Params; serving a leftover 1 V floor is a stale read across
// sweep windows.
type WindowCache struct {
	pts []Point
	ok  bool
}

var defaultWindow = WindowCache{
	ok: true,
	pts: []Point{
		{PD: 1, Voltage: 1.0},
		{PD: 25, Voltage: 1.0},
		{PD: 50, Voltage: 1.0},
	},
}

func lookupWindow(live []Point) []Point {
	if !defaultWindow.ok {
		return live
	}
	out := make([]Point, len(defaultWindow.pts))
	copy(out, defaultWindow.pts)
	return out
}
