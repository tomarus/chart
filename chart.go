// Package chart generates interactive svg or png charts from time series data.
package chart

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/tomarus/chart/axis"
	"github.com/tomarus/chart/data"
	"github.com/tomarus/chart/image"
	"github.com/tomarus/chart/palette"
)

// Chart is the main chart type used for all operations.
type Chart struct {
	width, height    int
	marginx, marginy int
	start, end       int64
	title            string
	data             data.Collection
	palette          *palette.Palette
	image            image.Image
	writer           io.Writer
	axes             []*axis.Axis
}

// Options defines a type used to initialize a Chart using NewChart()
type Options struct {
	Title         string      // guess what, leave empty to hide
	Size          string      // big is 1440px, small is 720px, auto is size of dataset
	Width, Height int         // overrides Size
	Scheme        string      // palette colorscheme, default "white"
	Theme         string      // if random scheme is used, set to "light" to use light colors, otherwise a dark theme is generated
	Start, End    int64       // start + end epoch of data
	Image         image.Image // the chart image type, chart.SVG{} or chart.PNG{}
	W             io.Writer   // output writer to write image to
	Axes          []*axis.Axis
}

// Render renders the final image to the io.Writer.
func (c *Chart) Render() error {
	if len(c.data) == 0 {
		return fmt.Errorf("no data available")
	}
	if c.width < 100 {
		return fmt.Errorf("image too small, set size or width or supply more datapoints")
	}
	c.data.Normalize(c.height)
	sort.Sort(c.data)

	for i := range c.data {
		c.data[i].Scale = c.axes[1].Scales(c.height, 0, c.data[i].Max)
	}

	c.image.Start(c.writer, c.width, c.height, c.marginx, c.marginy, c.start, c.end, c.palette, c.data)

	err := c.image.Graph()
	if err != nil {
		return err
	}

	c.axes[0].Draw(c.image, c.width, c.height, c.marginx, c.marginy, float64(c.start), float64(c.end))
	c.axes[1].Draw(c.image, c.width, c.height, c.marginx, c.marginy, 0, c.data[0].Max)

	c.drawTitle(c.width+c.marginx, c.height)
	c.image.Legend()
	c.image.Border(c.marginx-1, c.marginy-1, c.width+1, c.height+1)
	return c.image.End()
}

// drawTitle sets the chart title.
func (c *Chart) drawTitle(width, height int) {
	if c.title == "" {
		return
	}
	c.image.Text("title", "right", image.TitleRole, width-4, 12+2, c.title)
}

// AddData adds a single data set.
func (c *Chart) AddData(opt *data.Options, d []float64) (err error) {
	if opt.Type == "" {
		opt.Type = "area"
	}
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
	}
	c.data = append(c.data, newdata)
	return err
}

// NewChart initializes a new svg chart.
func NewChart(o *Options) (*Chart, error) {
	w := o.W
	if w == nil {
		w = os.Stdout
	}
	c := &Chart{title: o.Title, marginx: 48, marginy: 20, image: o.Image, writer: w, data: data.Collection{}, axes: o.Axes}

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

	c.start = o.Start
	c.end = o.End

	if o.Scheme != "" {
		c.palette, _ = palette.NewPalette(o.Scheme, o.Theme)
	} else {
		c.palette, _ = palette.NewPalette("white")
	}
	return c, nil
}
