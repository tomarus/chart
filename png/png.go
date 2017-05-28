package png

import (
	"fmt"
	"image"
	"image/draw"
	pngo "image/png"
	"io"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/tomarus/chart/data"
	"github.com/tomarus/chart/palette"
)

// PNG implements the chart interface to write PNG images.
type PNG struct {
	w                io.Writer
	width, height    int
	marginx, marginy int
	start, end       int
	pal              *palette.Palette
	img              *image.Paletted
}

func New() *PNG {
	return &PNG{}
}

func (png *PNG) Start(wr io.Writer, w, h, mx, my, start, end int, p *palette.Palette) {
	png.w = wr
	png.width = w
	png.height = h
	png.marginx = mx
	png.marginy = my
	png.start = start
	png.end = end
	png.pal = p
}

func (png *PNG) End() error {
	return pngo.Encode(png.w, png.img)
}

func (png *PNG) Graph(d data.Collection) {
	png.img = image.NewPaletted(image.Rect(0, 0, png.width+png.marginx+4, png.height+(2*png.marginy)+(d.Len()*16)), png.pal.Palette)

	bg := image.NewUniform(png.pal.GetColor("background"))
	draw.Draw(png.img, png.img.Bounds(), bg, image.ZP, draw.Src)

	for pt := range d {
		data := d[pt]
		col := png.pal.GetAxisColorName(pt)
		a := float64(data.NMax) / float64(png.height)
		b := float64(data.NMax) - a*float64(png.height)
		for i := range data.Values {
			v := int(float64(data.Values[i])*a + b)
			png.Line(col, i+png.marginx, png.height+png.marginy, i+png.marginx, png.height-v+png.marginy)
		}
	}
}

func (png *PNG) Text(col, align string, x, y int, txt string) {
	fill := image.NewUniform(png.pal.GetColor("title")) // Palette[2])
	d := &font.Drawer{
		Dst:  png.img,
		Src:  fill,
		Face: basicfont.Face7x13,
		Dot: fixed.Point26_6{
			X: fixed.I(x),
			Y: fixed.I(y),
		},
	}

	switch align {
	case "middle", "center":
		d.Dot.X -= d.MeasureString(txt) / 2
	case "end", "right":
		d.Dot.X -= d.MeasureString(txt)
	}
	d.DrawString(txt)
}

func (png *PNG) ID(id string) {
}

func (png *PNG) EndID() {
}

func (png *PNG) Line(color string, x1, y1, x2, y2 int) {
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	skip := 1
	if color == "grid" {
		skip = 2
	}
	ruler := png.pal.GetColor(color)
	if x1 == x2 {
		for i := y1; i < y2; i += skip {
			png.img.Set(x1, i, ruler)
		}
	} else if y1 == y2 {
		for i := x1; i < x2; i += skip {
			png.img.Set(i, y1, ruler)
		}
	}
}

func (png *PNG) rectFill(color string, x1, y1, w, h int) {
	for i := 0; i < h; i++ {
		png.Line(color, x1, y1+i, x1+w, y1+i)
	}
}

func (png *PNG) Legend(d data.Collection, p *palette.Palette) {
	x := png.marginx
	y := png.height + png.marginy + 4

	maxstrlen := 0
	for i := range d {
		if len(d[i].Title) >= maxstrlen {
			maxstrlen = len(d[i].Title)
		}
	}
	for i := range d {
		png.rectFill(p.GetAxisColorName(i), x, y+16+(i*16), 12, 12)

		min, max, avg := d[i].MinMaxAvg()
		mmax := data.FormatSI(max, 1, 1000, "", "", "")
		mmin := data.FormatSI(min, 1, 1000, "", "", "")
		mavg := data.FormatSI(avg, 1, 1000, "", "", "")
		q := fmt.Sprintf("%%-%ds    Max: %%6s    Avg: %%6s    Min: %%6s", maxstrlen)
		s := fmt.Sprintf(q, d[i].Title, mmax, mavg, mmin)
		png.Text("title", "left", x+20, y+26+(i*16), s)
	}
}

func (png *PNG) Border(x, y, w, h int) {
	c := "border"
	png.Line(c, x, y, x+w, y)
	png.Line(c, x+w, y, x+w, y+h)
	png.Line(c, x+w, y+h, x, y+h)
	png.Line(c, x, y+h, x, y)
}
