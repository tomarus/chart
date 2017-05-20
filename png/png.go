package png

import (
	"image"
	"image/draw"
	"io"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"

	"github.com/tomarus/chart/data"
	"github.com/tomarus/chart/palette"

	pngo "image/png"
)

type PNG struct {
	w                io.Writer
	width, height    int
	marginx, marginy int
	start, end       int
	pal              *palette.Palette
	img              *image.Paletted
}

func New(w io.Writer) (*PNG, error) {
	return &PNG{w: w}, nil
}

func (png *PNG) Start(w, h, mx, my, start, end int, p *palette.Palette) {
	png.width = w
	png.height = h
	png.marginx = mx
	png.marginy = my
	png.start = start
	png.end = end
	png.pal = p

	png.img = image.NewPaletted(image.Rect(0, 0, w+mx+32, h+(2*my)), p.Palette)

	bg := image.NewUniform(p.GetColor("background"))
	draw.Draw(png.img, png.img.Bounds(), bg, image.ZP, draw.Src)
}

func (png *PNG) End() error {
	return pngo.Encode(png.w, png.img)
}

func (png *PNG) Graph(d data.Collection) {
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

func (png *PNG) Text(color, align string, x, y int, txt string) {
	size := 12.
	dpi := 72.

	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Println(err)
		return
	}

	fill := image.NewUniform(png.pal.GetColor(color))
	d := &font.Drawer{
		Dst: png.img,
		Src: fill,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: font.HintingFull,
		}),
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
	if color == "gridlines" {
		color = "grid"
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

func (png *PNG) Legend(data.Collection, *palette.Palette) {
	// TODO
}

func (png *PNG) Border(x, y, w, h int) {
	c := "border"
	png.Line(c, x, y, x+w, y)
	png.Line(c, x+w, y, x+w, y+h)
	png.Line(c, x+w, y+h, x, y+h)
	png.Line(c, x, y+h, x, y)
}
