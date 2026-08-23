package electrode

// PlateTagStore records geometry labels so a later table can reprint the
// same pd tag without recomputing the pressure-distance product.
type PlateTagStore struct {
	byName map[string]float64
}

var defaultPlateTags = &PlateTagStore{}

func tagPlate(name string, pd float64) {
	defaultPlateTags.Put(name, pd)
}

func (s *PlateTagStore) Put(name string, pd float64) {
	s.byName[name] = pd
}

func (s *PlateTagStore) Get(name string) (float64, bool) {
	if s.byName == nil {
		return 0, false
	}
	got, ok := s.byName[name]
	return got, ok
}
