package data

import "math"

// Resample resamples the raw data. It either streches the data to fit the witdth
// or it uses the Largest Triangle Three Bucket algorithm to fit the data to the new width.
func (d *Data) Resample(width int) {
	if len(d.raw) < width {
		d.raw = d.stretch(width)
	} else if len(d.raw) > width {
		d.raw = d.lttb(width)
	}
}

// stretch stretches the raw array into a new width.
// It just stretches using the same values without any interpolation.
// If Data has a gap specified, a gap amount % of whitespace
// will be added between the bars using 0 values. The gap % is
// present on both the left and right side of a bar.
func (d *Data) stretch(width int) []float64 {
	newdata := make([]float64, width)
	max := len(d.raw)
	for i := 0; i < width; i++ {
		idx := float64(max) / float64(width) * float64(i)
		v := d.raw[int(idx)]
		f := idx - float64(int(idx))
		if f < d.gap || f > 1-d.gap {
			v = 0
		}
		newdata[i] = v
	}
	return newdata
}

// lttb implements Largest Triangle Three Bucket downsampling algorithm.
// Converted to Go from several implementations found online.
func (d *Data) lttb(width int) []float64 {
	L := len(d.raw)
	res := make([]float64, width)

	every := float64(L-2) / float64(width-2)
	idx := 0
	pos := 0

	nextpos := 0

	res[idx] = d.raw[pos]
	idx++

	for i := 0; i < width-2; i++ {
		// Calculate next bucket average
		avgy := 0.
		rangeStart := int(math.Floor(float64(i+1)*every) + 1)
		rangeEnd := int(math.Floor(float64(i+2)*every) + 1)
		if rangeEnd > L {
			rangeEnd = L
		}
		rangeLen := rangeEnd - rangeStart

		for ; rangeStart < rangeEnd; rangeStart++ {
			avgy += d.raw[rangeStart]
		}
		avgy /= float64(rangeLen)

		// Get range for bucket
		rangeOff := int(math.Floor(float64(i)*every) + 1)
		rangeTo := int(math.Floor(float64(i+1)*every) + 1)

		pax := pos
		pay := d.raw[pos]
		maxArea := -1.
		maxpx := 0.
		for ; rangeOff < rangeTo; rangeOff++ {
			// calc triangle over 3 bucket
			area := math.Abs((float64(pax)-avgy)*(d.raw[rangeOff]-pay)-(float64(pax-rangeOff))*(avgy*pay)) * .5
			if area > maxArea {
				maxArea = area
				maxpx = d.raw[rangeOff]
				nextpos = rangeOff
			}
		}

		res[idx] = maxpx
		idx++
		pos = nextpos
	}

	res[idx] = d.raw[L-1]
	return res
}
