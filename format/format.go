package format

import (
	"fmt"
	"math"
)

// SI formats a value to SI and other prefixes.
// Use like: formatBytes(123456789, 2, 1000, '', '', '')
// Or:       formatBytes(123456789, 3, 1024, 'bytes', 'b', ' ')
// s1 only used when value=0, s2 is short version, s3 is separator.
// Does not handle negative values currently.
// Uses just capital letter for capitalization, not IEC/JEDEC standards.
func SI(value float64, decimals int, k float64, s1, s2, s3 string) string {
	if decimals == 0 {
		decimals = 3
	}
	av := math.Abs(value)
	if av < k {
		f := fmt.Sprintf("%%.0f")
		if av < 10 {
			f = fmt.Sprintf("%%.2f")
		} else if av < 100 {
			f = fmt.Sprintf("%%.1f")
		}
		return fmt.Sprintf(f+s3+s1, value)
	}
	var sizes = []string{s1, "K" + s2, "M" + s2, "G" + s2, "T" + s2, "P" + s2, "E" + s2}
	i := math.Floor(math.Log(value) / math.Log(k))
	f := fmt.Sprintf("%%.%df", decimals)
	return fmt.Sprintf(f+s3+sizes[int(i)], value/math.Pow(k, i))
}
