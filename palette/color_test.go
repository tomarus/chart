package palette

import (
	"image/color"
	"testing"
)

func TestConvert(t *testing.T) {
	x := col2hex(&color.RGBA{128, 192, 64, 255})
	if x != "#80c040" {
		t.Fatal(x)
	}
	x = col2hex(&color.RGBA{80, 32, 11, 255})
	if x != "#50200b" {
		t.Fatal(x)
	}
}

func TestParseColor(t *testing.T) {
	var testParse = []struct {
		in  string
		out color.RGBA
	}{
		{"#fff", color.RGBA{255, 255, 255, 255}},
		{"#123", color.RGBA{17, 34, 51, 255}},
		{"#f0c0a0", color.RGBA{240, 192, 160, 255}},
		{"#01f203", color.RGBA{1, 242, 3, 255}},
		{"#01f20380", color.RGBA{1, 242, 3, 128}},
	}

	for _, tp := range testParse {
		res, err := ParseColor(tp.in)
		if err != nil {
			t.Fatal(err)
		}
		if *res != tp.out {
			t.Errorf("%s does not match %v, got %v", tp.in, tp.out, res)
		}
	}

	_, err := ParseColor("white")
	if err == nil {
		t.Error("Expected error parsing \"white\"")
	}
}
