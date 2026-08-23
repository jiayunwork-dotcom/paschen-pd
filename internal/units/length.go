package units

// LengthFromMm converts a length in millimetres to centimetres.
func LengthFromMm(mm float64) float64 {
	return mm / 10.0
}

// LengthToMm converts a length in centimetres to millimetres.
func LengthToMm(cm float64) float64 {
	return cm * 10.0
}

// LengthFromM converts a length in metres to centimetres.
func LengthFromM(m float64) float64 {
	return m * 100.0
}

// LengthToM converts a length in centimetres to metres.
func LengthToM(cm float64) float64 {
	return cm / 100.0
}
