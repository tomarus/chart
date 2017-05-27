package chart

import (
	"bufio"
	"bytes"
	"testing"
	"time"

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
		Start:  uint64(time.Now().AddDate(0, 0, -1).Unix()),
		End:    uint64(time.Now().Unix()),
		Xdiv:   12,
		Ydiv:   5,
		W:      w,
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
	if c.palette.GetHexColor("background") != "#fff" {
		t.Fatal("default scheme should be white")
	}

	// test data

	c.AddData("area", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
	if c.width != 12 {
		// A width of 12px is actually unviewable
		t.Errorf("Expected width of 10, got %d", c.width)
	}
	c.Render()

	c, _ = NewChart(opts)
	err = c.Render()
	if err == nil {
		t.Fatal("expected error no data available")
	}

	c, _ = NewChart(opts)
	c.AddData("area", []float64{1, 2, 3, 4, 5, 6})
	err = c.Render()
	if err == nil {
		t.Fatal("expected error xdiv <= datalen")
	}

	opts.Width = 720
	opts.Height = 540
	c, _ = NewChart(opts)
	c.AddData("area", []float64{1, 2, 3, 4, 5, 6})
	err = c.Render()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// t.Log(out.String())
}
