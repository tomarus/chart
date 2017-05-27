package colors

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

func TestColorOneEight(t *testing.T) {
	test := color.RGBA{128, 160, 80, 255}
	c := NewColor1(&test).RGBA()
	if *c != test {
		t.Error("Color1 != color.RGBA")
	}
}
