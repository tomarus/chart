package colors

import "image/color"

func OneTo8(in float64) uint8 {
	return uint8(in * 255.)
}

func EightTo1(in uint8) float64 {
	return float64(in) * 1. / 255.
}

type Color1 struct {
	R, G, B, A float64
}

func NewColor1(c *color.RGBA) *Color1 {
	return &Color1{EightTo1(c.R), EightTo1(c.G), EightTo1(c.B), EightTo1(c.A)}
}

func (c *Color1) RGBA() *color.RGBA {
	return &color.RGBA{OneTo8(c.R), OneTo8(c.G), OneTo8(c.B), OneTo8(c.A)}
}
