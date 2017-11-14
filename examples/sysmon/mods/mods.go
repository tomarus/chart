package mods

const maxKeep = 900

// Dataset is a set of datapoints and a description.
type Dataset struct {
	Title  string
	Values []float64
}

// Collector is the interface used to descrbie monitoring/updater methods.
type Collector interface {
	Len() int
	Update() error
	Data() []Dataset
}

// store is a fixed length slice
type store struct {
	values []float64
}

func (s *store) set(v float64) {
	if len(s.values) >= maxKeep {
		s.values = s.values[1:maxKeep]
	}
	s.values = append(s.values, v)
}
