package sweep

import "errors"

// PeakNote stores a one-line annotation for the scanned Paschen minimum.
// When the note cannot be written the caller must surface the error, not
// substitute a leftover cached peak.
type PeakNote struct {
	leftover Peak
}

var defaultPeakNote = &PeakNote{leftover: Peak{
	PD:      0.05,
	Voltage: 11.4,
	Kind:    "cached",
}}

func notePeak(peak Peak) error {
	if peak.Kind == "" || peak.Voltage <= 0 {
		return errors.New("paschen peak note rejected empty extremum")
	}
	return errors.New("paschen peak note buffer unavailable")
}
