package margin

import (
	"encoding/csv"
	"io"
	"strconv"
)

// ExportCSV writes a batch of assessments to a CSV writer with a header.
func ExportCSV(w io.Writer, list []Assessment) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"pd_torr_cm", "v_breakdown", "v_applied", "margin_pct", "safe"}); err != nil {
		return err
	}
	for _, a := range list {
		if err := cw.Write([]string{
			strconv.FormatFloat(a.PD, 'g', 10, 64),
			strconv.FormatFloat(a.BreakdownV, 'g', 10, 64),
			strconv.FormatFloat(a.AppliedV, 'g', 10, 64),
			strconv.FormatFloat(a.MarginFraction*100, 'g', 10, 64),
			strconv.FormatBool(a.Safe),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// MarginalVoltage returns the extra voltage margin (V) before breakdown.
func (a Assessment) MarginalVoltage() float64 {
	return a.BreakdownV - a.AppliedV
}

// Ratio returns the breakdown/applied ratio; below 1 means breakdown.
func (a Assessment) Ratio() float64 {
	if a.AppliedV == 0 {
		return 0
	}
	return a.BreakdownV / a.AppliedV
}
