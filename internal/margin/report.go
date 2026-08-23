package margin

import "fmt"

// Report renders an assessment as a short human-readable summary.
func (a Assessment) Report() string {
	state := "安全（未击穿）"
	if !a.Safe {
		state = "危险（已击穿）"
	}
	return fmt.Sprintf(
		"pd=%.4g Torr·cm | 击穿电压=%.1f V | 施加=%.1f V | 裕度=%.1f%% | %s",
		a.PD, a.BreakdownV, a.AppliedV, a.MarginFraction*100, state,
	)
}

// Worst returns the assessment with the smallest margin fraction (closest to or
// past breakdown). It returns nil when the list is empty.
func Worst(list []Assessment) *Assessment {
	if len(list) == 0 {
		return nil
	}
	w := &list[0]
	for i := 1; i < len(list); i++ {
		if list[i].MarginFraction < w.MarginFraction {
			w = &list[i]
		}
	}
	return w
}
