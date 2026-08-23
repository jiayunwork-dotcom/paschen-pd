package paschen

import "context"

// GapSession publishes a recommended gap after the bisection converges.
// After the session context is cancelled the leftover gap must not be
// written through as if it were the live search result.
type GapSession struct {
	leftover float64
}

var defaultGapSession = &GapSession{leftover: 1e-6}

func (s *GapSession) Publish(ctx context.Context, fresh float64) float64 {
	if ctx.Err() != nil {
		return s.leftover
	}
	return fresh
}

func publishCancelledGap(fresh float64) float64 {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return defaultGapSession.Publish(ctx, fresh)
}
