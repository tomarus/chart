package palette

import (
	"image/color"
	"testing"
)

func TestConv(t *testing.T) {
	u8 := OneTo8(0.5)
	t.Log(u8)

	f64 := EightTo1(128)
	t.Log(f64)
}

func TestHSV(t *testing.T) {
	var td = []struct {
		H HSV
		E string
	}{
		{HSV{0, 1, 1, 1}, "#ff0000"},
		{HSV{90, 1, 1, 1}, "#7fff00"},
		{HSV{180, 1, 1, 1}, "#00ffff"},
		{HSV{270, 1, 1, 1}, "#7f00ff"},
		{HSV{360, 1, 1, 1}, "#ff0000"},
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

func TestColorOneEight(t *testing.T) {
	test := color.RGBA{128, 160, 80, 255}
	c := NewColor1(&test).RGBA()
	if *c != test {
		t.Error("Color1 != color.RGBA")
	}
}

func TestString(t *testing.T) {
	test := NewHSVA(360, .5, .5, 1)
	s := test.String()
	if s != "h360.00 s0.50 v0.50 a1.00" {
		t.Errorf("String() did not match (%s)", s)
	}
}
