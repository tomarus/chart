// Package chart generates interactive svg charts from time series data.
//
// Example
//
// import "github.com/tomarus/chart"
//
// opts := &chart.Options{
// 	Title:  "Traffic",
//	Image:  svg.New(), // or png.New()
// 	Size:   "big",     // big is 1440px, small is 720px, auto is size of dataset
//	Height: 300,       // Defaults to -1, when size=auto height is set to width/4, otherwise set fixed height
//	Width:  900,       // If a width is supplied, height is implied and both are used in stead of size setting
// 	Scheme: "white",   // or black/random/pink/solarized or hsl:180,0.5,0.25
// 	Start:  start_epoch,
// 	End:    end_epoch,
// 	Xdiv:   12,
// 	Ydiv:   5,
// 	W:      w,
// }
//
// c, err := chart.NewChart(opts)
// if err != nil {
// 	panic(err)
// }
//
// warn := c.AddData("area", []yourData)
// if err != nil {
// 	fmt.Println(warn)
// }
//
// w.Header().Set("Content-Type", "image/svg+xml")
// c.Render()
//
package chart

import (
	"fmt"
	"io"
	"sort"

	"github.com/tomarus/chart/data"
	"github.com/tomarus/chart/palette"
)

// Image defines the interface for image (svg/png) backends.
type Image interface {
	Start(wr io.Writer, w, h, mx, my, start, end int, p *palette.Palette)
	End() error
	Graph(data.Collection)
	Text(color, align string, x, y int, txt string)
	ID(id string)
	EndID()
	Line(color string, x1, y1, x2, y2 int)
	Legend(data.Collection, *palette.Palette)
	Border(x, y, w, h int)
}

// Chart is the main chart type used for all operations.
type Chart struct {
	width, height    int
	marginx, marginy int
	start, end       float64
	xdiv, ydiv       int
	title            string
	data             data.Collection
	palette          *palette.Palette
	image            Image
	writer           io.Writer
}

// Options defines a type used to initialize a Chart using NewChart()
type Options struct {
	Title         string    // guess what, leave empty to hide
	Size          string    // big is 1440px, small is 720px, auto is size of dataset
	Width, Height int       // overrides Size
	Scheme        string    // palette colorscheme, default "white"
	Theme         string    // if random scheme is used, set to "light" to use light colors, otherwise a dark theme is generated
	Start, End    uint64    // start + end epoch of data
	Xdiv, Ydiv    int       // num grid divisions (default x12 y5)
	Image         Image     // the chart image type, chart.SVG{} or chart.PNG{}
	W             io.Writer // output writer to write image to
}

// Render renders the final image to the io.Writer.
func (c *Chart) Render() error {
	c.data.Normalize(c.height)
	sort.Sort(c.data)
	c.scales(c.ydiv)
	c.image.Start(c.writer, c.width, c.height, c.marginx, c.marginy, int(c.start), int(c.end), c.palette)
	c.image.Graph(c.data)
	c.xgrid(c.marginx, c.marginy, c.xdiv, c.start, c.end)
	c.ygrid(c.marginx, c.marginy, c.ydiv)
	c.drawTitle(c.width+c.marginx, c.height)
	c.image.Legend(c.data, c.palette)
	c.image.Border(c.marginx-1, c.marginy-1, c.width+1, c.height+1)
	return c.image.End()
}

// drawTitle sets the chart title.
func (c *Chart) drawTitle(width, height int) {
	if c.title == "" {
		return
	}
	c.image.Text("title", "right", width-4, 12+2, c.title)
}

// AddData adds a single datasource of type t and data d
// A warning is returned if the graph + data sizes do not match.
func (c *Chart) AddData(t string, d []float64) (err error) {
	newdata := data.NewData(t, d)
	if len(d) == 0 {
		c.data = append(c.data, newdata)
		return fmt.Errorf("Added empty dataset")
	}

	// Setup auto width if not done so already.
	if c.width == -1 {
		c.width = len(d)
		if c.height == -1 {
			c.height = c.width / 4
		}
	}

	if len(d) != c.width {
		newdata.Resample(c.width)
		err = fmt.Errorf("Resampling data from %d to %d", len(d), c.width)
	}

	c.data = append(c.data, newdata)
	return err
}

// NewChart initializes a new svg chart.
func NewChart(o *Options) (*Chart, error) {
	c := &Chart{title: o.Title, marginx: 48, marginy: 20, image: o.Image, writer: o.W}

	if o.Size == "big" {
		c.width = 1440
		c.height = 360
	} else if o.Size == "small" {
		c.width = 720
		c.height = 240
	} else { // "auto"
		c.width = -1
		c.height = -1
	}
	if o.Width > 0 {
		c.width = o.Width
	}
	if o.Height > 0 {
		c.height = o.Height
	}

	c.xdiv = 12
	c.ydiv = 5
	if o.Xdiv > 0 {
		c.xdiv = o.Xdiv
	}
	if o.Ydiv > 0 {
		c.ydiv = o.Ydiv
	}

	c.start = float64(o.Start * 1000.0)
	c.end = float64(o.End * 1000.0)

	if o.Scheme != "" {
		c.palette = palette.NewPalette(o.Scheme, o.Theme)
	} else {
		c.palette = palette.NewPalette("white")
	}
	return c, nil
}
