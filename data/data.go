package data

// Data contains a single set of data most likely imported from tsm.
type Data struct {
	raw    []float64
	Max    float64  `json:"fmax"`   // max raw value
	NMax   int      `json:"max"`    // max normalized value
	Scale  []string `json:"scale"`  // yaxis labels
	Values []int    `json:"values"` // pixel values
	Type   string   `json:"type"`
	Title  string   `json:"title"`
}

// Options contains configuration for a single dataset.
type Options struct {
	Type  string
	Title string
}

// NewData creates a new dataset from []float64.
func NewData(opt *Options, in []float64) Data {
	return Data{Type: opt.Type, Title: opt.Title, raw: in}
}

// Len returns the number of items in the dataset.
func (d *Data) Len() int {
	return len(d.raw)
}

// MinMaxAvg returns the Minimum, Maximum and Average values of the raw data.
func (d *Data) MinMaxAvg() (float64, float64, float64) {
	max := 0.
	avg := 0.
	min := 0.
	for _, v := range d.raw {
		if max < v {
			max = v
		}
		if v != 0 && (min == 0 || min > v) {
			min = v
		}
		avg += v
	}
	avg /= float64(len(d.raw))
	d.Max = max
	return min, max, avg
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
