package palette

import (
	"fmt"
	"image/color"
	"math/rand"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/tomarus/chart/colors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var axisColors = []string{"area", "color1", "color2", "color3"}

// regexps used to match a randomized scheme based on a single hsl value.
var mx1 = regexp.MustCompile("hsl:([0-9.]+),([0-9.]+),([0-9.]+)")
var mx2 = regexp.MustCompile("hsl:([0-9a-f]{6})")

// Palette contains information for working with palettes and colors.
type Palette struct {
	sync.RWMutex
	palette map[string]color.Color
}

// NewPalette returns a color.Palette from the specified color scheme.
func NewPalette(opts ...string) (*Palette, error) {
	p := &Palette{}

	scheme := "white"
	if len(opts) > 0 {
		scheme = opts[0]
	}
	light := false
	if len(opts) > 1 && opts[1] == "light" {
		light = true
	}

	if scheme == "random" {
		p.randomize(true, 0, 1, 1, light)
	} else if mx1.MatchString(scheme) {
		x := mx1.FindStringSubmatch(scheme)
		fx, _ := strconv.ParseFloat(x[1], 10)
		fsat, _ := strconv.ParseFloat(x[2], 10)
		fval, _ := strconv.ParseFloat(x[3], 10)
		p.randomize(false, fx, fsat, fval, light)
	} else if mx2.MatchString(scheme) {
		x := mx2.FindStringSubmatch(scheme)
		ch, _ := colors.NewHSLHex(x[1])
		p.randomize(false, ch.H, ch.S, ch.L, light)
	} else if defs, x := getDefaultScheme(scheme); x {
		err := p.installPalette(defs)
		if err != nil {
			return nil, err
		}
	} else {
		defs, _ := getDefaultScheme("white")
		p.installPalette(defs)
	}
	return p, nil
}

// GetColor retrieves a color.Color from the palette (e.g. "background", "grid")
func (p *Palette) GetColor(name string) color.Color {
	p.RLock()
	defer p.RUnlock()
	if c, ok := p.palette[name]; ok {
		return c
	}
	return nil // &UnknownColor
}

// GetHexColor retrieves a hexadecimal color from the palette (e.g. "background", "grid")
func (p *Palette) GetHexColor(name string) string {
	return col2hex(p.GetColor(name))
}

// GetAxisColorName gets the color for the Nth datapoint.
func (p *Palette) GetAxisColorName(id int) string {
	return axisColors[id]
}

// GetHexAxisColor gets the color for the Nth datapoint.
func (p *Palette) GetHexAxisColor(id int) string {
	return p.GetHexColor(axisColors[id])
}

func (p *Palette) installPalette(pal map[string]string) error {
	p.palette = map[string]color.Color{}
	for k, c := range pal {
		col, err := ParseColor(c)
		if err != nil {
			return fmt.Errorf("installPalette: parsecolor: %v", err)
		}
		p.palette[k] = col
	}
	return nil
}

func (p *Palette) randomize(r bool, hue, sat, val float64, light bool) {
	if r {
		hue = rand.Float64() * 360
		sat = 1. / 3. // rand.Float64()/3 + .25
		val = .5
	}

	bg := colors.NewHSL(hue, 0.15*sat, 0.1*val).RGBA()
	fg := colors.NewHSL(hue, 0.05*sat, 2.*val).RGBA()
	grid := colors.NewHSLA(hue, 0, .2, 1).RGBA()
	border := colors.NewHSLA(hue, 0, .5, 1).RGBA()
	marker := colors.NewHSL(hue, sat, 1.5*val).RGBA()

	// defaults to dark theme
	if light {
		fg = bg
		bg = colors.NewHSL(hue, 0.15*sat, 1).RGBA()
		grid = colors.NewHSLA(hue, 0, 0, .25).RGBA()
		marker = colors.NewHSL(hue, sat, 0.5*val).RGBA()
	}

	p.Lock()
	defer p.Unlock()
	p.palette = map[string]color.Color{
		"background": bg,
		"title":      fg,
		"grid":       grid,
		"border":     border,
		"marker":     marker,
		"select":     marker,
		"area":       colors.NewHSL(hue, sat, val).RGBA(),
		"color1":     colors.NewHSL(hue, 0.75*sat, 1.25*val).RGBA(),
		"color2":     colors.NewHSL(hue, 0.5*sat, 1.5*val).RGBA(),
		"color3":     colors.NewHSL(hue, 0.25*sat, 0.75*val).RGBA(),
	}
}
