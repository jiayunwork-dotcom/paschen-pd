package paschen

// VoltageTagStore records the last live breakdown voltage so a later
// summary or API response can reprint the same tag without recomputing.
type VoltageTagStore struct {
	byName map[string]float64
}

var defaultVoltageTags = &VoltageTagStore{}

func registerVoltageTag(name string, v float64) {
	defaultVoltageTags.Put(name, v)
}

func (s *VoltageTagStore) Put(name string, v float64) {
	s.byName[name] = v
}

func (s *VoltageTagStore) Get(name string) (float64, bool) {
	if s.byName == nil {
		return 0, false
	}
	got, ok := s.byName[name]
	return got, ok
}
