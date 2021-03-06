package colors

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func col2hex(in *color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", in.R, in.G, in.B)
}
func TestHSL(t *testing.T) {
	var td = []struct {
		H HSL
		E string
	}{
		{HSL{0, 1, .25, 1}, "#7f0000"},
		{HSL{90, 1, .25, 1}, "#3f7f00"},
		{HSL{180, 1, .25, 1}, "#007f7f"},
		{HSL{270, 1, .25, 1}, "#3f007f"},
		{HSL{360, 1, .25, 1}, "#7f0000"},
		//
		{HSL{0, .5, .75, 1}, "#df9f9f"},
		{HSL{90, .5, .75, 1}, "#bfdf9f"},
		{HSL{180, .5, .75, 1}, "#9fdfdf"},
		{HSL{270, .5, .75, 1}, "#bf9fdf"},
		{HSL{360, .5, .75, 1}, "#df9f9f"},
		//
		{HSL{0, .5, .5, 1}, "#bf3f3f"},
		{HSL{90, .5, .5, 1}, "#7fbf3f"},
		{HSL{180, .5, .5, 1}, "#3fbfbf"},
		{HSL{270, .5, .5, 1}, "#7f3fbf"},
		{HSL{360, .5, .5, 1}, "#bf3f3f"},
	}

	for _, test := range td {
		x := col2hex(test.H.RGBA())
		if x != test.E {
			t.Errorf("Expected %s got %s", test.E, x)
		}
	}
}

func cmp(a, b float64) bool {
	const e = 1e-2
	return (a-b) < e && (b-a) < e
}

func TestHSLHex(t *testing.T) {
	c, err := NewHSLHex("invalid")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	c, err = NewHSLHex("8040c0")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp(c.H, 180.) || !cmp(c.S, 0.25) || !cmp(c.L, 0.75) {
		t.Errorf("Expected h=180, s=.25, l=.75 got %v", c)
	}
}

func TestHSLGrey(t *testing.T) {
	r := NewHSL(180, 0., .5).RGBA()
	if r.R != r.G || r.G != r.B || r.B != r.R {
		t.Error("Should have greyscale color")
	}
}

func TestHSLOverflow(t *testing.T) {
	c := NewHSL(400, .5, .5)
	if !cmp(c.H, 40) {
		t.Errorf("Hue should be 40, got %v", c.H)
	}
	c = NewHSL(-40, .5, .5)
	if !cmp(c.H, 320) {
		t.Errorf("Hue should be 320, got %v", c.H)
	}
}

func TestHSLString(t *testing.T) {
	test := NewHSLA(360, .5, .5, 1)
	s := test.String()
	if s != "h360.00 s0.50 l0.50 a1.00" {
		t.Errorf("String() did not match (%s)", s)
	}
}

func TestHSLImage(t *testing.T) {
	t.Skip()
	f, err := os.OpenFile("/tmp/output-hsl.png", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 360, 200))
	for h := 0; h < 360; h++ {
		for v := 0; v < 200; v++ {
			col := NewHSL(float64(h), 1., 1./200.*float64(v))
			img.Set(h, v, col.RGBA())
		}
	}
	png.Encode(f, img)
}
