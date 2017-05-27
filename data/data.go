package data

// Data contains a single set of data most likely imported from tsm.
type Data struct {
	raw    []float64
	Max    float64  `json:"fmax"`   // max raw value
	NMax   int      `json:"max"`    // max normalized value
	Scale  []string `json:"scale"`  // yaxis labels
	Values []int    `json:"values"` // pixel values
	Type   string   `json:"type"`
}

// NewData creates a new dataset from []float64.
func NewData(typ string, in []float64) Data {
	return Data{Type: typ, raw: in}
}

// Len returns the number of items in the dataset.
func (d *Data) Len() int {
	return len(d.raw)
}

// Normalize normalizes the raw/tsm values to height.
func (d *Data) normalize(height int) {
	d.Max = 0.
	for _, v := range d.raw {
		if d.Max < v {
			d.Max = v
		}
	}

	if d.Max == 0 {
		// we have an empty dataset
		for _, v := range d.raw {
			d.Values = append(d.Values, int(v))
		}
		return
	}

	fmax := float64(height)
	a := fmax / d.Max
	b := fmax - a*d.Max

	for _, v := range d.raw {
		newv := a*v + b
		d.Values = append(d.Values, int(newv))
	}
}

// normalizeMax normalizes our max value according to height and a global max value.
func (d *Data) normalizeMax(height int, max float64) {
	fmax := float64(height)
	a2 := fmax / max
	b2 := fmax - a2*max
	d.NMax = int(a2*float64(d.Max) + b2)
}
