package paschen

// Physical and numeric constants used throughout the Paschen model.
//
// The Paschen coefficients A, B and gamma are empirical; the constants below
// are the unit normalisation factors that let callers convert between the
// common pressure/length conventions without losing precision.
const (
	// TorrPerPa is the number of Torr in one Pascal (1 Torr = 133.322 Pa).
	TorrPerPa = 1.0 / 133.322368421

	// PaPerTorr is the number of Pascal in one Torr.
	PaPerTorr = 133.322368421

	// MmPerCm is millimetres per centimetre.
	MmPerCm = 10.0

	// CmPerM is centimetres per metre.
	CmPerM = 100.0

	// AtmPerTorr is standard atmospheres per Torr.
	AtmPerTorr = 1.0 / 760.0

	// BarPerTorr is bar per Torr.
	BarPerTorr = 1.0 / 750.0616827
)

// StandardGravity is the gravitational acceleration used as a fallback when a
// caller does not supply one (m/s^2). It is exported so applications can stay
// consistent with the model's defaults.
const StandardGravity = 9.80665

// MachineEpsilon is a small value used to guard divisions and comparisons near
// the Paschen cutoff.
const MachineEpsilon = 1e-12
