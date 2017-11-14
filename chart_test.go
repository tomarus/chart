package chart

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/tomarus/chart/axis"
	"github.com/tomarus/chart/data"
	"github.com/tomarus/chart/image"
	"github.com/tomarus/chart/png"
	"github.com/tomarus/chart/svg"
)

func TestChart(t *testing.T) {
	var out bytes.Buffer
	w := bufio.NewWriter(&out)

	opts := &Options{
		Title:  "Test Title",
		Image:  svg.New(),
		Size:   "small",
		Scheme: "random",
		Theme:  "light",
		Start:  time.Now().AddDate(0, 0, -1).Unix(),
		End:    time.Now().Unix(),
		Axes: []*axis.Axis{
			axis.NewTime(axis.Bottom, "01-02 15:04").Ticks(12),
			axis.NewSI(axis.Left, 1000).Ticks(5),
		},
		W: w,
	}

	// test sizes

	c, err := NewChart(opts)
	if err != nil {
		t.Fatal(err)
	}
	if c.width != 720 {
		t.Fatal("width should be 720")
	}

	opts.Size = "big"
	c, _ = NewChart(opts)
	if c.width != 1440 {
		t.Fatal("width should be 1440")
	}

	opts.Width = 320
	opts.Height = 240
	c, _ = NewChart(opts)
	if c.width != 320 || c.height != 240 {
		t.Fatal("expected width/ehgith 320x240")
	}

	// test palette

	opts.Width = 0
	opts.Height = 0
	opts.Size = "auto"
	opts.Scheme = ""
	c, _ = NewChart(opts)
	if c.palette.GetHexColor("background") != "#ffffff" {
		t.Fatal("default scheme should be white")
	}

	// test data

	c.AddData(&data.Options{Title: "testing"}, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
	if c.width != 12 {
		// A width of 12px is actually unviewable
		t.Errorf("Expected width of 10, got %d", c.width)
	}

	c, _ = NewChart(opts)
	err = c.Render()
	if err == nil {
		t.Fatal("expected error no data available")
	}

	c, _ = NewChart(opts)
	c.AddData(&data.Options{Title: "testing"}, []float64{1, 2, 3, 4, 5, 6})
	err = c.Render()
	if err == nil {
		t.Fatal("expected error xdiv <= datalen")
	}

	opts.Width = 720
	opts.Height = 540
	c, _ = NewChart(opts)
	c.AddData(&data.Options{}, []float64{1, 2, 3, 4, 5, 6})
	err = c.Render()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	// TODO: actually test svg output somehow

	opts.Image = png.New()
	c, _ = NewChart(opts)
	c.AddData(&data.Options{}, []float64{1, 2, 3, 4, 5, 6})
	err = c.Render()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	// TODO: actually test png output somehow
}

func testimg(img image.Image) {
	var out bytes.Buffer
	w := bufio.NewWriter(&out)

	opts := &Options{
		Title:  "Test Title",
		Image:  img,
		Size:   "small",
		Scheme: "white",
		Theme:  "light",
		Start:  time.Now().AddDate(0, 0, -1).Unix(),
		End:    time.Now().Unix(),
		Axes: []*axis.Axis{
			axis.NewTime(axis.Bottom, "01-02 15:04").Duration(4 * time.Hour).Grid(4),
			axis.NewSI(axis.Left, 1000).Ticks(4).Grid(2),
		},
		W: w,
	}

	c, _ := NewChart(opts)
	c.AddData(&data.Options{}, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
	c.Render()
}

func BenchmarkSVG(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testimg(svg.New())
	}
}

func BenchmarkPNG(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testimg(png.New())
	}
}

func ExampleChart() {
	opts := &Options{
		Title:  "Title on top of the chart",
		Image:  svg.New(), // or png.New()
		Size:   "big",     // big is 1440px, small is 720px, auto is size of dataset
		Height: 300,       // Defaults to -1, when size=auto height is set to width/4, otherwise set fixed height
		Width:  900,       // If a width is supplied, height is implied and both are used in stead of size setting
		Scheme: "white",   // or black/random/pink/solarized or hsl:180,0.5,0.25
		Theme:  "light",   // default is dark.
		Start:  time.Now().AddDate(0, 0, -1).Unix(),
		End:    time.Now().Unix(),
		W:      os.Stdout,
		Axes: []*axis.Axis{
			axis.NewTime(axis.Bottom, "01-02 15:04").Duration(4 * time.Hour).Grid(4),
			axis.NewSI(axis.Left, 1000).Ticks(4).Grid(2),
			//
			// - Example custom time format
			// axis.New(axis.Bottom, func(in float64) string {
			//	return time.Unix(int64(in), 0).Format("01-02 15:04")
			// 	return in.(time.Time).Format("2006-01-02")
			// }).Duration(4 * time.Hour).Grid(),
			//
			// - Example other custom format
			// axis.New(axis.Left, func(in float64) string {
			// 	return fmt.Sprintf("%.1f", in/3.14)
			// }).Ticks(5).Grid(),
		},
	}

	c, err := NewChart(opts)
	if err != nil {
		panic(err)
	}

	exdata := make([]float64, 256)
	for i := 0; i < 255; i++ {
		exdata[i] = float64(i)
	}
	warn := c.AddData(&data.Options{Type: "area", Title: "My Data Description"}, exdata)
	if err != nil {
		fmt.Println(warn)
	}
	c.Render()
}
