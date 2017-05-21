package palette

import (
	"fmt"
	"image/color"
	"math"
)

type HSV struct {
	H, S, V, A float64
}

func (hsv *HSV) String() string {
	return fmt.Sprintf("h%.2f s%.2f v%.2f a%.2f", hsv.H, hsv.S, hsv.V, hsv.A)
}

func NewHSV(h, s, v float64) *HSV {
	return NewHSVA(h, s, v, 1.)
}

func NewHSVA(h, s, v, a float64) *HSV {
	if h > 360 {
		h -= 360
	} else if h < 0 {
		h += 360
	}
	return &HSV{h, s, v, a}
}

func (hsv *HSV) RGBA() *color.RGBA {
	if hsv.S == 0 {
		c := &Color1{hsv.V, hsv.V, hsv.V, hsv.A}
		return c.RGBA()
	}
	c := &Color1{0, 0, 0, hsv.A}
	h := hsv.H / 60
	if h == 6 {
		h = 0
	}
	i := math.Floor(h)
	v1 := hsv.V * (1 - hsv.S)
	v2 := hsv.V * (1 - hsv.S*(h-i))
	v3 := hsv.V * (1 - hsv.S*(1-(h-i)))
	switch int(i) {
	case 0:
		c.R = hsv.V
		c.G = v3
		c.B = v1
	case 1:
		c.R = v2
		c.G = hsv.V
		c.B = v1
	case 2:
		c.R = v1
		c.G = hsv.V
		c.B = v3
	case 3:
		c.R = v1
		c.G = v2
		c.B = hsv.V
	case 4:
		c.R = v3
		c.G = v1
		c.B = hsv.V
	default:
		c.R = hsv.V
		c.G = v1
		c.B = v2
	}
	return c.RGBA()
}
