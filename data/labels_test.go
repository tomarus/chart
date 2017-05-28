package data

import "testing"

func TestYLabels(t *testing.T) {
	d := NewData(&Options{"line", ""}, []float64{0, 5})
	d.normalize(5)
	d.CreateScale(5)
	expect := []string{"5", "4", "3", "2", "1"}
	if !seq(d.Scale, expect) {
		t.Errorf("Expected %#v got %#v", expect, d.Scale)
	}
}
