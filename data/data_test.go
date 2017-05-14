package data

import (
	"sort"
	"testing"
)

var testData = Collection{
	NewData("line", []float64{1, 2, 3, 4, 5}),
	NewData("area", []float64{10, 20, 30, 40, 50}),
}

func eq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func seq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestNormalize(t *testing.T) {
	testData[0].normalize(10)
	expect := []int{2, 4, 6, 8, 10}
	if !eq(testData[0].Values, expect) {
		t.Errorf("Expected %#v got %#v", expect, testData[0].Values)
	}
}

func TestNormalizeMax(t *testing.T) {
	testData[0].normalizeMax(1000, 20)
	if testData[0].NMax != 250 {
		t.Errorf("NMax should be 250, is %d", testData[0].NMax)
	}
}

func TestChartNormalize(t *testing.T) {
	testData.Normalize(10)
	if testData[0].Max != 5 {
		t.Error("max should be 5")
	}
	if testData[1].Max != 50 {
		t.Error("max should be 50")
	}
}

func TestSort(t *testing.T) {
	// normalize should be called first to init Max
	sort.Sort(testData)
	if testData[0].Type != "area" {
		t.Errorf("Data is not sorted correctly (%#v)", testData)
	}
}
