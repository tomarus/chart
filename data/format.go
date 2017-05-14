package data

import (
	"fmt"
	"math"
)

// formatSI formats a value to SI and other prefixes.
// Use like: formatBytes(123456789, 2, 1000, '', '', '')
// Or:       formatBytes(123456789, 3, 1024, 'bytes', 'b', ' ')
// s1 only used when value=0, s2 is short version, s3 is separator.
// Does not handle negative values currently.
// Uses just capital letter for capitalization, not IEC/JEDEC standards.
func formatSI(value float64, decimals int, k float64, s1, s2, s3 string) string {
	if value == 0 {
		return "0" + s3 + s1
	}
	if decimals == 0 {
		decimals = 3
	}
	var sizes = []string{s1, "K" + s2, "M" + s2, "G" + s2, "T" + s2, "P" + s2, "E" + s2}
	i := math.Floor(math.Log(value) / math.Log(k))
	if math.IsNaN(i) || i < 0 {
		return "0" + s3 + s1
	}
	f := fmt.Sprintf("%%.%df", decimals)
	if i == 0 {
		f = fmt.Sprintf("%%.0f")
	}
	return fmt.Sprintf(f+s3+sizes[int(i)], value/math.Pow(k, i))
}
