// Package chart generates interactive svg or png charts from time series data.
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
	// Start initializes a new image and sets the defaults.
	Start(wr io.Writer, w, h, mx, my int, start, end int64, p *palette.Palette)

	// End finishes and writes the image to the output writer.
	End() error

	// Graph renders all chart dataset values to the visible chart area.
	Graph(data.Collection)

	// Text writes a string to the image.
	Text(color, align string, x, y int, txt string)

	// Line draws a line between two points. The current implementations need to
	// only draw horizontal or vertical lines.
	Line(color string, x1, y1, x2, y2 int)

	// Legend draws the image specific legend.
	Legend(data.Collection, *palette.Palette)

	// Border draws a border around the chart area.
	Border(x, y, w, h int)

	// ID is used to mark the start of the grid so svg/js can manipulate this.
	ID(id string)

	// EndID marks the end of the ID block specified earlier.
	EndID()
}

// Chart is the main chart type used for all operations.
type Chart struct {
	width, height    int
	marginx, marginy int
	start, end       int64
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
	Start, End    int64     // start + end epoch of data
	Xdiv, Ydiv    int       // num grid divisions (default x12 y5)
	Image         Image     // the chart image type, chart.SVG{} or chart.PNG{}
	W             io.Writer // output writer to write image to
}

// Render renders the final image to the io.Writer.
func (c *Chart) Render() error {
	if len(c.data) == 0 {
		return fmt.Errorf("no data available")
	}
	if c.xdiv >= c.data[0].Len() {
		return fmt.Errorf("xdivisions higher than dataset length")
	}
	c.data.Normalize(c.height)
	sort.Sort(c.data)
	c.scales(c.ydiv)
	c.image.Start(c.writer, c.width, c.height, c.marginx, c.marginy, c.start, c.end, c.palette)
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
func (c *Chart) AddData(opt *data.Options, d []float64) (err error) {
	newdata := data.NewData(opt, d)
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
	c := &Chart{title: o.Title, marginx: 48, marginy: 20, image: o.Image, writer: o.W, data: data.Collection{}}

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

	c.start = o.Start
	c.end = o.End

	if o.Scheme != "" {
		c.palette, _ = palette.NewPalette(o.Scheme, o.Theme)
	} else {
		c.palette, _ = palette.NewPalette("white")
	}
	return c, nil
}
