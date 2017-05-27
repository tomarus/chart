package colors

import (
	"image"
	"image/png"
	"os"
	"testing"
)

func TestHSV(t *testing.T) {
	var td = []struct {
		H HSV
		E string
	}{
		{HSV{0, 1, 1, 1}, "#ff0000"},
		{HSV{60, 1, 1, 1}, "#ffff00"},
		{HSV{120, 1, 1, 1}, "#00ff00"},
		{HSV{180, 1, 1, 1}, "#00ffff"},
		{HSV{240, 1, 1, 1}, "#0000ff"},
		{HSV{300, 1, 1, 1}, "#ff00ff"},
		//
		{HSV{0, .5, 1, 1}, "#ff7f7f"},
		{HSV{90, .5, 1, 1}, "#bfff7f"},
		{HSV{180, .5, 1, 1}, "#7fffff"},
		{HSV{270, .5, 1, 1}, "#bf7fff"},
		{HSV{360, .5, 1, 1}, "#ff7f7f"},
		//
		{HSV{0, .5, .5, 1}, "#7f3f3f"},
		{HSV{90, .5, .5, 1}, "#5f7f3f"},
		{HSV{180, .5, .5, 1}, "#3f7f7f"},
		{HSV{270, .5, .5, 1}, "#5f3f7f"},
		{HSV{360, .5, .5, 1}, "#7f3f3f"},
	}

	for _, test := range td {
		x := col2hex(test.H.RGBA())
		if x != test.E {
			t.Errorf("Expected %s got %s", test.E, x)
		}
	}
}

func TestHSVGrey(t *testing.T) {
	r := NewHSV(180, 0., .5).RGBA()
	if r.R != r.G || r.G != r.B || r.B != r.R {
		t.Error("Should have greyscale color")
	}
}

func TestHSVOverflow(t *testing.T) {
	c := NewHSV(400, .5, .5)
	if !cmp(c.H, 40) {
		t.Errorf("Hue should be 40, got %v", c.H)
	}
	c = NewHSV(-40, .5, .5)
	if !cmp(c.H, 320) {
		t.Errorf("Hue should be 320, got %v", c.H)
	}
}

func TestHSVString(t *testing.T) {
	test := NewHSVA(360, .5, .5, 1)
	s := test.String()
	if s != "h360.00 s0.50 v0.50 a1.00" {
		t.Errorf("String() did not match (%s)", s)
	}
}

func TestHSVImage(t *testing.T) {
	t.Skip()
	f, err := os.OpenFile("/tmp/output-hsv.png", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 360, 200))
	for h := 0; h < 360; h++ {
		for v := 0; v < 200; v++ {
			col := NewHSV(float64(h), 1., 1./200.*float64(v))
			img.Set(h, v, col.RGBA())
		}
	}
	png.Encode(f, img)
}
