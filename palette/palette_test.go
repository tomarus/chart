package palette

import "testing"

func TestDefaultPalette(t *testing.T) {
	pal, _ := NewPalette("") // defaults to "white"
	if len(pal.palette) != 10 {
		t.Fatal("len(pal) != 10")
	}
	if col2hex(pal.GetColor("background")) != "#ffffff" {
		t.Error("Expected white")
	}
}

func TestPalette(t *testing.T) {
	pal, _ := NewPalette("black")
	if len(pal.palette) != 10 {
		t.Fatal("len(pal) != 10")
	}

	h := pal.GetHexColor("grid")
	if h != "#404040" {
		t.Errorf("Unexpected hex grid color (%v) should be #404040", h)
	}

	c := pal.GetColor("grid")
	if col2hex(c) != "#404040" {
		t.Errorf("Unexpected grid color (%v) %s", c, col2hex(c))
	}

	c = pal.GetColor("doesnotexist")
	if c != nil {
		t.Error("Should be nil")
	}

	ch := pal.GetHexAxisColor(1)
	if ch != "#4444ff" {
		t.Error("Axis 1 color should be #00f")
	}

	n := pal.GetAxisColorName(0)
	if n != "area" {
		t.Error("First axis color should be area")
	}
}

func TestRandom(t *testing.T) {
	p, _ := NewPalette("random")
	if len(p.palette) != 10 {
		t.Errorf("Palette should have 10 entries.")
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
