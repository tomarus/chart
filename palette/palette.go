package palette

import (
	"image/color"
	"math/rand"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var palettes = map[string]map[string]string{
	"white": {
		"background": "#fff",
		"title":      "#660",
		"grid":       "#999",
		"border":     "#ccc",
		"marker":     "#666",
		"select":     "#999",
		"area":       "#8abb31",
		"color1":     "#0000ff",
		"color2":     "#ff0000",
		"color3":     "#ff00ff",
	},
	"black": {
		"background": "#000",
		"title":      "#fff",
		"grid":       "#ccc",
		"border":     "#ccc",
		"marker":     "#fff",
		"select":     "#0f0",
		"area":       "#8abb31",
		"color1":     "#4444ff",
		"color2":     "#ff4444",
		"color3":     "#ff44ff",
	},
	"pink": {
		"background": "#f2d1ba",
		"title":      "#5e2728",
		"grid":       "#5e2728",
		"border":     "#5e2728",
		"marker":     "#7e4748",
		"select":     "#eee",
		"area":       "#f34093",
		"color1":     "#f78bd1",
		"color2":     "#d2082d",
		"color3":     "#cf0063",
	},
	"solarized": { // http://ethanschoonover.com/solarized
		"background": "#002b36",
		"title":      "#eee8d5",
		"grid":       "#586e75",
		"border":     "#657b83",
		"marker":     "#fdf6e3",
		"select":     "#fdf6e3",
		"area":       "#93a1a1",
		"color1":     "#657b83",
		"color2":     "#586e75",
		"color3":     "#073642",
	},
}

var colorOrder = []string{"background", "title", "grid", "border", "marker", "select", "area", "color1", "color2", "color3"}
var axisColors = []string{"area", "color1", "color2", "color3"}

// UnknownColor describes an unknown color to the palette. Currently this is bright purple.
var UnknownColor = color.RGBA{255, 0, 255, 255}

var lock sync.RWMutex

// Palette contains information for working with palettes and colors.
type Palette struct {
	Name    string
	Palette []color.Color
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func mkrandom(key string, r bool, hue, sat, val float64) {
	if r {
		hue = rand.Float64() * 360
		sat = rand.Float64()/2 + .375
	}
	lock.Lock()
	defer lock.Unlock()
	palettes[key] = map[string]string{
		"background": col2hex(NewHSV(hue, sat/3, val/6).RGBA()),
		"title":      col2hex(NewHSV(hue, 0, val).RGBA()),
		"grid":       col2hex(NewHSV(hue, 0, val/2).RGBA()),
		"border":     col2hex(NewHSV(hue, 0, val/2).RGBA()),
		"marker":     col2hex(NewHSV(hue, sat/4, val/1.5).RGBA()),
		"select":     col2hex(NewHSV(hue, sat/3, val/2).RGBA()),
		"area":       col2hex(NewHSV(hue, sat, val/4*3).RGBA()),
		"color1":     col2hex(NewHSV(hue+120, sat, val/4*3).RGBA()),
		"color2":     col2hex(NewHSV(hue+240, sat, val/4*3).RGBA()),
		"color3":     col2hex(NewHSV(hue+120, sat, val/4*3).RGBA()),
	}

}

var mx1 = regexp.MustCompile("hsv:([0-9.]+),([0-9.]+),([0-9.]+)")
var mx2 = regexp.MustCompile("hsv:([0-9a-f]{6})")

// NewPalette returns a color.Palette from the specified color scheme.
func NewPalette(opts ...string) *Palette {
	scheme := "white"
	if len(opts) > 0 {
		scheme = opts[0]
	}

	if scheme == "random" {
		mkrandom("random", true, 0, 1, 1)
	}
	if mx1.MatchString(scheme) {
		x := mx1.FindStringSubmatch(scheme)
		fx, _ := strconv.ParseFloat(x[1], 10)
		fsat, _ := strconv.ParseFloat(x[2], 10)
		fval, _ := strconv.ParseFloat(x[3], 10)
		mkrandom(scheme, false, fx, fsat, fval)
	}
	if mx2.MatchString(scheme) {
		x := mx2.FindStringSubmatch(scheme)
		ch, _ := hex2hsv(x[1])
		mkrandom(scheme, false, ch.H, ch.S, ch.V)
	}

	lock.RLock()
	if _, x := palettes[scheme]; !x {
		scheme = "white"
	}
	lock.RUnlock()

	pal := make([]color.Color, len(colorOrder))
	for i, col := range colorOrder {
		lock.RLock()
		pal[i], _ = ParseColor(palettes[scheme][col])
		lock.RUnlock()
	}
	return &Palette{Name: scheme, Palette: pal}
}

// GetColor retrieves a color.Color from the palette (e.g. "background", "grid")
func (p *Palette) GetColor(name string) color.Color {
	lock.RLock()
	defer lock.RUnlock()
	for i := range colorOrder {
		if colorOrder[i] == name {
			return p.Palette[i]
		}
	}
	return UnknownColor
}

// GetHexColor retrieves a hexadecimal color from the palette (e.g. "background", "grid")
func (p *Palette) GetHexColor(name string) string {
	lock.RLock()
	defer lock.RUnlock()
	return palettes[p.Name][name]
}

// GetAxisColorName gets the color for the Nth datapoint.
func (p *Palette) GetAxisColorName(id int) string {
	return axisColors[id]
}

// GetHexAxisColor gets the color for the Nth datapoint.
func (p *Palette) GetHexAxisColor(id int) string {
	lock.RLock()
	defer lock.RUnlock()
	return palettes[p.Name][axisColors[id]]
}
