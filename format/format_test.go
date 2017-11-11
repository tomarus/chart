package format

import (
	"testing"
)

var testFormat = []struct {
	v          float64
	d          int
	k          float64
	s1, s2, s3 string
	expect     string
}{
	{0, 2, 1000, "xx", "x", "", "0.00xx"},
	{123, 2, 1000, "bytes", "b", " ", "123 bytes"},
	{123, 2, 1000, "", "b", "", "123"},
	{-1, 0, 1000, "", "", "", "-1.00"},
	{25, 0, 1000, "", "", "", "25.0"},
	{-25, 0, 1000, "", "", "", "-25.0"},
	{-.1, 1, 1000, "", "", "", "-0.10"},
	{.1, 1, 1000, "", "", "", "0.10"},
	{5, 1, 1000, "", "", "", "5.00"},
	{1234.5, 0, 1000, "bytes", "b", " ", "1.234 Kb"},
	{1234.5, 2, 1000, "bytes", "b", " ", "1.23 Kb"},
	{1234.5, 2, 1024, "bytes", "b", " ", "1.21 Kb"},
	{1234456.7, 2, 1024, "apples", "a", " ", "1.18 Ma"}, // mega apples
}

func TestFormatSI(t *testing.T) {
	for _, f := range testFormat {
		res := SI(f.v, f.d, f.k, f.s1, f.s2, f.s3)
		if res != f.expect {
			t.Errorf("Expected %s got %s testdata %v", f.expect, res, f)
		}
	}
}
