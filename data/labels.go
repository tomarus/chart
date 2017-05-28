package data

// CreateScale creates an array of Scale labels for y axis.
func (d *Data) CreateScale(n int) {
	d.Scale = []string{}
	max := d.Max
	skippy := max / float64(n)
	for iy := 0; iy < n; iy++ {
		d.Scale = append(d.Scale, FormatSI(max, 1, 1000, "", "", ""))
		max -= skippy
	}
}
