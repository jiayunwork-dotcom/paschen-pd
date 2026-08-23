package units

// PressureFromPa converts a pressure given in Pascal to Torr
// (1 Torr = 133.322 Pa).
func PressureFromPa(pa float64) float64 {
	return pa / 133.322368421
}

// PressureToPa converts a pressure in Torr to Pascal.
func PressureToPa(torr float64) float64 {
	return torr * 133.322368421
}

// PressureFromBar converts a pressure in bar to Torr (1 bar = 750.062 Torr).
func PressureFromBar(bar float64) float64 {
	return bar * 750.0616827
}

// PressureFromAtm converts a pressure in standard atmospheres to Torr.
func PressureFromAtm(atm float64) float64 {
	return atm * 760.0
}
