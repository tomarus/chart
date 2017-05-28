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
		data.NewData(&data.Options{Type: "line"}, []float64{0, 5}),
		data.NewData(&data.Options{Type: "line"}, []float64{0, 1, 2, 3, 4, 5}),
	},
}

func TestXLabels(t *testing.T) {
	t0, _ := time.Parse("2006-01-02", "2017-01-01")
	t1, _ := time.Parse("2006-01-02", "2017-01-30")
	labelChart.start = float64(t0.Unix() * 1000)
	labelChart.end = float64(t1.Unix() * 1000)
	labelChart.width = 1440

	labels := labelChart.xlabels(12)
	for i, lbl := range labels {
		exp := time.Unix(t0.Unix()+(int64(i+1)*((29*86400)/12)), 0).Format("01-02 15:04")
		if lbl != exp {
			t.Errorf("Expected %v got %v", exp, lbl)
		}
	}
}
