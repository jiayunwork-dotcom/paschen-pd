package sweep

import (
	"encoding/csv"
	"io"
	"strconv"

	"paschen-pd/internal/paschen"
)

// WriteCSV writes the sampled curve as CSV with headers "pd_torr_cm" and "v_v".
func WriteCSV(w io.Writer, c Config) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"pd_torr_cm", "v_v"}); err != nil {
		return err
	}
	for _, p := range Run(c) {
		if err := cw.Write([]string{
			strconv.FormatFloat(p.PD, 'g', 10, 64),
			strconv.FormatFloat(p.Voltage, 'g', 10, 64),
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// SummaryLine returns a one-line characterization of the curve for logging.
func SummaryLine(c Config) string {
	pdMin, vMin := Stats(c)
	return "min V=" + strconv.FormatFloat(vMin, 'g', 6, 64) +
		" V at pd=" + strconv.FormatFloat(pdMin, 'g', 6, 64) +
		" Torr·cm"
}

// WithAir returns a sweep config using the dry-air coefficients.
func WithAir(from, to float64, points int, log bool) Config {
	return Config{
		From:   from,
		To:     to,
		Points: points,
		Log:    log,
		Params: paschen.DefaultAir(),
	}
}
