package api

// ValidationError describes why a request could not be processed.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// ValidateBreakdown checks a breakdown request for physical sanity. It returns a
// ValidationError when p or d is non-positive, otherwise nil.
func ValidateBreakdown(req BreakdownRequest) error {
	if req.P <= 0 {
		return ValidationError{Field: "p", Message: "pressure must be positive"}
	}
	if req.D <= 0 {
		return ValidationError{Field: "d", Message: "gap distance must be positive"}
	}
	return nil
}

// ValidateSweep checks a sweep request. It additionally requires d_max to exceed
// d_min so the scanned interval is non-empty.
func ValidateSweep(req SweepRequest) error {
	if req.P <= 0 {
		return ValidationError{Field: "p", Message: "pressure must be positive"}
	}
	if req.DMin <= 0 {
		return ValidationError{Field: "d_min", Message: "d_min must be positive"}
	}
	if req.DMax <= req.DMin {
		return ValidationError{Field: "d_max", Message: "d_max must exceed d_min"}
	}
	return nil
}

// IsValidation reports whether err is a ValidationError.
func IsValidation(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(ValidationError)
	return ok
}
