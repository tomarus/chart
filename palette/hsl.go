package palette

import (
	"fmt"
	"image/color"
)

// Based on https://github.com/THEjoezack/ColorMine

type HSL struct {
	H, S, L, A float64
}

func (hsl *HSL) String() string {
	return fmt.Sprintf("h%.2f s%.2f l%.2f a%.2f", hsl.H, hsl.S, hsl.L, hsl.A)
}

func NewHSL(h, s, l float64) *HSL {
	return NewHSLA(h, s, l, 1.)
}

func NewHSLA(h, s, l, a float64) *HSL {
	if h > 360 {
		h -= 360
	}
	return &HSL{h, s, l, a}
}

func (hsl *HSL) RGBA() *color.RGBA {
	if hsl.S == 0 {
		c := &Color1{hsl.L, hsl.L, hsl.L, hsl.A}
		return c.RGBA()
	}

	t2 := 0.
	if hsl.L < 0.5 {
		t2 = hsl.L * (1.0 + hsl.S)
	} else {
		t2 = hsl.L + hsl.S - (hsl.L * hsl.S)
	}
	t1 := 2.*hsl.L - t2

	rng := hsl.H / 360
	r := hsl.gc(t1, t2, rng+1./3.)
	g := hsl.gc(t1, t2, rng)
	b := hsl.gc(t1, t2, rng-1./3.)
	c := &Color1{r, g, b, hsl.A}
	return c.RGBA()
}

func (hsl *HSL) gc(t1, t2, t3 float64) float64 {
	t3 = hsl.rng(t3)
	if t3 < 1./6. {
		return t1 + (t2-t1)*6.*t3
	}
	if t3 < .5 {
		return t2
	}
	if t3 < 2./3. {
		return t1 + ((t2 - t1) * ((2. / 3.) - t3) * 6.)
	}
	return t1
}

func (hsl *HSL) rng(t float64) float64 {
	if t < 0. {
		return t + 1.
	}
	if t > 1. {
		return t - 1.
	}
	return t
}
