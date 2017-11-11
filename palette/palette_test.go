package palette

import "testing"

func TestDefaultPalette(t *testing.T) {
	pal, _ := NewPalette("") // defaults to "white"
	if len(pal.palette) != 12 {
		t.Fatal("len(pal) != 12")
	}
	if col2hex(pal.GetColor("background")) != "#ffffff" {
		t.Error("Expected white")
	}
}

func TestPalette(t *testing.T) {
	pal, _ := NewPalette("black")
	if len(pal.palette) != 12 {
		t.Fatal("len(pal) != 12")
	}

	h := pal.GetHexColor("grid")
	if h != "#999999" {
		t.Errorf("Unexpected hex grid color (%v) should be #999999", h)
	}

	c := pal.GetColor("grid")
	if col2hex(c) != "#999999" {
		t.Errorf("Unexpected grid color (%v) %s", c, col2hex(c))
	}

	c = pal.GetColor("doesnotexist")
	if c != unknownColor {
		t.Error("Should be unknownColor")
	}

	ch := pal.GetHexAxisColor(1)
	if ch != "#295135" {
		t.Error("Axis 1 color should be #295135")
	}

	n := pal.GetAxisColorName(0)
	if n != "area" {
		t.Error("First axis color should be area")
	}
}

func TestRandom(t *testing.T) {
	p, _ := NewPalette("random")
	if len(p.palette) != 12 {
		t.Errorf("Palette should have 12 entries.")
	}

	p, _ = NewPalette("hsl:180,0.5,0.5")
	c := col2hex(p.GetColor("area"))
	if c != "#3fbfbf" {
		t.Errorf("Area color should be #3dbfbf is %s", c)
	}

	p, _ = NewPalette("hsl:804020", "light")
	c = col2hex(p.GetColor("area"))
	if c != "#172828" {
		t.Errorf("Area color should be #172828 is %s", c)
	}
}
