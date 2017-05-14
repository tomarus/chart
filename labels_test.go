package chart

import (
	"testing"
	"time"

	"github.com/tomarus/chart/data"
)

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

var labelChart = Chart{
	data: data.Collection{
		data.NewData("line", []float64{0, 5}),
		data.NewData("line", []float64{0, 1, 2, 3, 4, 5}),
	},
}

func TestXLabels(t *testing.T) {
	t0, _ := time.Parse("2006-01-02", "2017-01-01")
	t1, _ := time.Parse("2006-01-02", "2017-01-30")
	labelChart.start = float64(t0.Unix() * 1000)
	labelChart.end = float64(t1.Unix() * 1000)
	labelChart.width = 1440
	// XXX variations possible but not tested!
	expect := []string{"01-03 11:00", "01-05 21:00", "01-08 07:00", "01-10 17:00", "01-13 03:00", "01-15 13:00", "01-17 23:00", "01-20 09:00", "01-22 19:00", "01-25 05:00", "01-27 15:00"}
	labels := labelChart.xlabels(12)
	if !seq(labels, expect) {
		t.Errorf("Expected %#v got %#v", expect, labels)
	}
}
