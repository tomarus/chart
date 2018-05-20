package png

import (
	"fmt"
	"io"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gosmallcapsitalic"

	"github.com/tomarus/chart/data"
	"github.com/tomarus/chart/format"
	myimg "github.com/tomarus/chart/image"
	"github.com/tomarus/chart/palette"
)

// PNG implements the chart interface to write PNG images.
type PNG struct {
	w                io.Writer
	gg               *gg.Context
	data             data.Collection
	width, height    int
	marginx, marginy int
	start, end       int64
	pal              *palette.Palette
}

// New initializes a new png chart image writer.
func New() *PNG {
	return &PNG{}
}

// Start initializes a new image and sets the defaults.
func (png *PNG) Start(wr io.Writer, w, h, mx, my int, start, end int64, p *palette.Palette, d data.Collection) {
	png.w = wr
	png.data = d
	png.width = w
	png.height = h
	png.marginx = mx
	png.marginy = my
	png.start = start
	png.end = end
	png.pal = p
}

// End finishes and writes the image to the output writer.
func (png *PNG) End() error {
	return png.gg.EncodePNG(png.w)
}

// Graph renders all chart dataset values to the visible chart area.
func (png *PNG) Graph() error {
	png.gg = gg.NewContext(png.width+png.marginx+4, png.height+(2*png.marginy)+((png.data.Len()+1)*16))
	png.gg.SetColor(png.pal.GetColor("background"))
	png.gg.Clear()

	for pt, data := range png.data {
		col := png.pal.GetAxisColorName(pt)
		a := float64(data.NMax) / float64(png.height)
		b := float64(data.NMax) - a*float64(png.height)
		for i := range data.Values {
			if data.Values[i] < 0 {
				return fmt.Errorf("Negative values not supported")
			}
			v := int(float64(data.Values[i])*a + b)
			png.Line(col, i+png.marginx, png.height+png.marginy, i+png.marginx, png.height-v+png.marginy)
		}
	}
	return nil
}

// face returns the font face to use. If the role is set to "title" a larger font is used.
func (png *PNG) face(role myimg.TextRole) {
	var ttfont *truetype.Font
	size := 13.
	dpi := 72.
	h := font.HintingNone

	if role == myimg.GridRole {
		f, err := truetype.Parse(gomono.TTF)
		if err != nil {
			panic(err)
		}
		ttfont = f
	} else {
		f, err := truetype.Parse(gosmallcapsitalic.TTF)
		if err != nil {
			panic(err)
		}
		ttfont = f
		size = 16.
	}

	face := truetype.NewFace(ttfont, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: h,
	})
	png.gg.SetFontFace(face)
}

// Text writes a string to the image.
func (png *PNG) Text(col, align string, role myimg.TextRole, x, y int, txt string) {
	ax := 0.
	switch align {
	case "middle", "center":
		ax = .5
	case "end", "right":
		ax = 1
	}
	png.gg.SetColor(png.pal.GetColor(col))
	png.face(role)
	png.gg.DrawStringAnchored(txt, float64(x), float64(y), ax, 0)
}

// TextID writes a string to the image.
func (png *PNG) TextID(id, col, align string, role myimg.TextRole, x, y int, txt string) {
	png.Text(col, align, role, x, y, txt)
}

// Line draws a line between the points using the color name from the palette.
func (png *PNG) Line(color string, x1, y1, x2, y2 int) {
	ruler := png.pal.GetColor(color)
	if color == "grid" || color == "grid2" {
		png.gg.SetDash(1)
		png.gg.SetLineWidth(.5)
		png.gg.DrawLine(float64(x1), float64(y1), float64(x2), float64(y2))
		png.gg.SetColor(png.pal.GetColor(color))
		png.gg.Stroke()
	} else {
		png.gg.SetLineWidth(1)
		png.gg.DrawLine(float64(x1), float64(y1), float64(x2), float64(y2))
		png.gg.SetColor(ruler)
		png.gg.Stroke()
	}
}

func (png *PNG) rectFill(color string, x1, y1, w, h int) {
	for i := 0; i < h; i++ {
		png.Line(color, x1, y1+i, x1+w, y1+i)
	}
}

// Legend draws the image specific legend.
func (png *PNG) Legend() {
	x := png.marginx
	y := png.height + png.marginy + 4

	q := "Min     Max     Avg"
	png.Text("title2", "right", myimg.GridRole, x+png.width, y+26, q)
	y += 16

	for i, d := range png.data {
		png.rectFill(png.pal.GetAxisColorName(i), x, y+16, 12, 12)

		min, max, avg := d.MinMaxAvg()
		// FIXME use axis formatters for this.
		mmax := format.SI(max, 1, 1000, "", "", "")
		mmin := format.SI(min, 1, 1000, "", "", "")
		mavg := format.SI(avg, 1, 1000, "", "", "")
		q := fmt.Sprintf("%6s  %6s  %6s", mmin, mmax, mavg)
		png.Text("title", "left", myimg.GridRole, x+20, y+26, d.Title)
		png.Text("title", "right", myimg.GridRole, x+png.width, y+26, q)
		png.Line("grid2", x, y+26+3, x+png.width, y+26+3)
		y += 16
	}
}

// Border draws a border around the chart area.
func (png *PNG) Border(x, y, w, h int) {
	c := "border"
	png.Line(c, x, y, x+w, y)
	png.Line(c, x+w, y, x+w, y+h)
	png.Line(c, x+w, y+h, x, y+h)
	png.Line(c, x, y+h, x, y)
}
