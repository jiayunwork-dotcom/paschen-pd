package paschen

import "math"

// FitResult holds coefficients recovered from measurements.
type FitResult struct {
	A          float64
	B          float64
	Gamma      float64
	ResidualRMSE float64
}

// FitGammaFromMargin recovers the secondary-emission coefficient gamma given
// fixed A and B and a single (pd, V) measurement. The Paschen denominator is
// ln(A·pd) - ln(ln(1+1/gamma)); rearranging one point gives
//
//	ln(ln(1+1/gamma)) = ln(A·pd) - B·pd/V =: D
//	ln(1+1/gamma)     = exp(D)
//	gamma = 1 / ( exp(exp(D)) - 1 )
//
// which is exact for a single measurement (the RMSE is therefore 0).
func FitGammaFromMargin(a, b, pd, v float64) FitResult {
	if v <= 0 || pd <= 0 {
		return FitResult{}
	}
	d := math.Log(a*pd) - b*pd/v
	if d <= 0 {
		return FitResult{}
	}
	gamma := 1 / (math.Exp(math.Exp(d)) - 1)
	return FitResult{A: a, B: b, Gamma: gamma, ResidualRMSE: 0}
}

// FitFromPoints performs a least-squares fit of A and B to several (pd, V)
// measurements while holding gamma fixed. It linearises ln(V) = ln(B) +
// ln(pd) - ln(A*pd - ln(ln(1+1/gamma))); the dominant linear part is
// ln(V) = ln(B) + ln(pd) - ln(A*pd), so we solve for ln(B) and A via a small
// normal-equation system on the linearised residual.
func FitFromPoints(pts []Point, gamma float64) FitResult {
	n := len(pts)
	if n < 2 {
		return FitResult{Gamma: gamma}
	}
	logGamma := math.Log(math.Log(1 + 1/gamma))
	var s0, s1, s2, sy, sxy float64
	for _, p := range pts {
		x := math.Log(p.PD)
		y := math.Log(p.V)
		s0 += 1
		s1 += x
		s2 += x * x
		sy += y
		sxy += x * y
	}
	// We fit y = c0 + c1*x  where c0 = ln(B) - logGamma, c1 = 1 (since ln(V) ≈
	// ln(B)+ln(pd) - logGamma is exact only when A*pd term negligible). We allow
	// c1 to deviate to absorb the A*pd term as a first-order correction.
	denom := float64(n)*s2 - s1*s1
	if denom == 0 {
		return FitResult{Gamma: gamma}
	}
	c1 := (float64(n)*sxy - s1*sy) / denom
	c0 := (sy - c1*s1) / float64(n)
	b := math.Exp(c0 + logGamma)
	a := math.Max(1e-6, c1) // A enters as multiplier of pd in the log term
	residual := 0.0
	for _, p := range pts {
		model := math.Log(a*p.PD) - logGamma
		pred := math.Exp(math.Log(b) + math.Log(p.PD) - model)
		residual += (pred - p.V) * (pred - p.V)
	}
	rmse := math.Sqrt(residual / float64(n))
	return FitResult{A: a, B: b, Gamma: gamma, ResidualRMSE: rmse}
}

// Point is a single (pd, voltage) measurement used for fitting.
type Point struct {
	PD float64
	V  float64
}
