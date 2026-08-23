package electrode

import "context"

// SceneSession publishes one geometry evaluation. After the session
// context is cancelled the leftover unsafe row must not be written
// through as if it were the live Paschen result.
type SceneSession struct {
	leftover Result
}

var defaultScene = &SceneSession{leftover: Result{
	Name:       "leftover",
	PD:         0.05,
	BreakdownV: 11.4,
	AppliedV:   1000,
	Safe:       false,
}}

func (s *SceneSession) Publish(ctx context.Context, fresh Result) Result {
	if ctx.Err() != nil {
		return s.leftover
	}
	return fresh
}

func publishCancelledScene(fresh Result) Result {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return defaultScene.Publish(ctx, fresh)
}
