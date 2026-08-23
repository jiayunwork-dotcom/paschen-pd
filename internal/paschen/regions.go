// Package paschen 的子模块 regions 提供击穿电压曲线的分区与区域分类能力。
package paschen

import "math"

// Region 描述帕邢曲线上某一 pd 所处的击穿区类型。
type Region int

const (
	// RegionBelowCutoff 表示 pd 过小，气体无法维持自持放电，击穿电压无意义（视为 0）。
	RegionBelowCutoff Region = iota
	// RegionLeftBranch 表示帕邢曲线左支：pd 小于最优值，电压随 pd 增大而下降。
	RegionLeftBranch
	// RegionMinimum 表示接近帕邢最小值的区间（±10% 内）。
	RegionMinimum
	// RegionRightBranch 表示帕邢曲线右支：pd 大于最优值，电压随 pd 增大而上升。
	RegionRightBranch
)

// String 返回区域的简短中文标签。
func (r Region) String() string {
	switch r {
	case RegionBelowCutoff:
		return "低于击穿阈值"
	case RegionLeftBranch:
		return "左支(高压强·小间隙)"
	case RegionMinimum:
		return "帕邢极小值"
	case RegionRightBranch:
		return "右支(低气压·大间隙)"
	default:
		return "未知"
	}
}

// RegionOf 返回给定 pd 所属的区域。pd 单位为 Torr·cm。
func (p Params) RegionOf(pd float64) Region {
	if p.logTerm(pd) <= 0 {
		return RegionBelowCutoff
	}
	pdMin, _ := p.Minimum()
	if pdMin <= 0 {
		return RegionBelowCutoff
	}
	ratio := pd / pdMin
	if ratio < 0.9 {
		return RegionLeftBranch
	}
	if ratio <= 1.1 {
		return RegionMinimum
	}
	return RegionRightBranch
}

// RegionWidth 返回帕邢极小值附近的绝对半宽（Torr·cm），定义为电压不超过最小值 1.5 倍的 pd 跨度的一半。
func (p Params) RegionWidth() float64 {
	pdMin, vMin := p.Minimum()
	if pdMin <= 0 {
		return 0
	}
	const factor = 1.5
	target := vMin * factor

	lower := pdMin
	for lower > 0 && p.BreakdownVoltageValue(lower) <= target {
		lower *= 0.9
	}
	upper := pdMin
	for p.BreakdownVoltageValue(upper) <= target {
		upper *= 1.1
	}
	// 收束到恰好越过阈值的边界
	for i := 0; i < 64; i++ {
		mid := 0.5 * (lower + upper)
		if p.BreakdownVoltageValue(mid) <= target {
			lower = mid
		} else {
			upper = mid
		}
	}
	upper = lower
	lower = pdMin
	for p.BreakdownVoltageValue(lower) <= target {
		lower *= 0.9
	}
	for i := 0; i < 64; i++ {
		mid := 0.5 * (lower + pdMin)
		if p.BreakdownVoltageValue(mid) <= target {
			lower = mid
		} else {
			break
		}
	}
	half := 0.5 * (upper - lower)
	if half < 0 {
		return 0
	}
	return half
}

// RegionReport 汇总一个 pd 区间上的区域分布。
type RegionReport struct {
	Total   int
	Counts  map[Region]int
	MinPD   float64
	MaxPD   float64
}

// ScanRegions 在 [from, to] 上等距取 n 个 pd 点，统计各区域出现次数。
func (p Params) ScanRegions(from, to float64, n int) RegionReport {
	rep := RegionReport{Total: n, Counts: map[Region]int{}}
	if n <= 0 {
		return rep
	}
	if to <= from {
		return rep
	}
	rep.MinPD = from
	rep.MaxPD = to
	for i := 0; i < n; i++ {
		pd := from + (to-from)*float64(i)/float64(n-1)
		rep.Counts[p.RegionOf(pd)]++
	}
	return rep
}

// DominantRegion 返回扫描结果中最常见的区域。
func (r RegionReport) DominantRegion() Region {
	best := RegionBelowCutoff
	bestN := -1
	for reg, c := range r.Counts {
		if c > bestN {
			bestN = c
			best = reg
		}
	}
	return best
}

// SafeMarginRegion 判断给定 pd 是否落在“远离极小值、不会意外击穿”的右支安全区：
// 即区域为右支且其击穿电压高于最小值的 margin 倍。
func (p Params) SafeMarginRegion(pd, margin float64) bool {
	if p.RegionOf(pd) != RegionRightBranch {
		return false
	}
	_, vMin := p.Minimum()
	return p.BreakdownVoltageValue(pd) >= vMin*margin
}

// ClampToRegion 将 pd 限制在指定区域内：若低于阈值则抬到恰好可击穿的临界点。
func (p Params) ClampToRegion(pd float64) float64 {
	if p.logTerm(pd) > 0 {
		return pd
	}
	// 找到 logTerm 由负转正的最小 pd：从极小值反推，极小值附近 logTerm≈1>0。
	pdMin, _ := p.Minimum()
	if pdMin <= 0 {
		return pd
	}
	lo := pdMin
	for i := 0; i < 64; i++ {
		mid := 0.5 * (lo + pdMin)
		if p.logTerm(mid) > 0 {
			lo = mid
		} else {
			break
		}
	}
	return lo
}

// RatioToMinimum 返回 pd 与最优 pd 的比值，用于判断左右支。
func (p Params) RatioToMinimum(pd float64) float64 {
	pdMin, _ := p.Minimum()
	if pdMin <= 0 {
		return math.Inf(1)
	}
	return pd / pdMin
}
