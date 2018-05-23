// Package chart generates interactive svg or png charts from time series data.
package chart

import (
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"time"

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
	sibase           int
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
	SIBase        int         // SI Base for auto axis calculation, default is 1000.
	Axes          []*axis.Axis
}

// formats define the separator length, the time format and the
// number of subgrids. This is used if Options.Axis is not specified.
// It aims to look nice globally but some formats might need to be
// added in the future.
var formats = []struct {
	days int
	dur  time.Duration
	fmt  string
	grid int
}{
	{0, 4 * time.Hour, "15:04", 4},
	{3, 8 * time.Hour, "15:04", 4},
	{7, 1 * 24 * time.Hour, "02 Jan", 4},
	{14, 2 * 24 * time.Hour, "02 Jan", 2},
	{30, 5 * 24 * time.Hour, "02-01", 5},
	{60, 7 * 24 * time.Hour, "02-01", 2},
	{180, 30 * 24 * time.Hour, "02-01", 4},
	{365, 60 * 24 * time.Hour, "02-01-06", 4},
	{365 * 2, 90 * 24 * time.Hour, "02-01-06", 3},
}

func (c *Chart) addAxes() {
	days := int(math.Round(float64(c.end)-float64(c.start)) / 86400.0)
	dur := 8 * time.Hour
	ffmt := "Mon 15:04"
	grid := 4
	for _, f := range formats {
		if days+1 >= f.days {
			dur = f.dur
			ffmt = f.fmt
			grid = f.grid
		}
	}
	c.axes = []*axis.Axis{
		axis.NewTime(axis.Bottom, ffmt).Duration(dur).Grid(grid),
		axis.NewSI(axis.Left, 1000).Ticks(4).Grid(2),
	}
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

	if len(c.axes) == 0 {
		c.addAxes()
	}

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
	c.image.Legend(float64(c.sibase))
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
	c := &Chart{title: o.Title, marginx: 48, marginy: 20, image: o.Image, writer: w, data: data.Collection{}, axes: o.Axes, sibase: o.SIBase}

	if c.sibase == 0 {
		c.sibase = 1000
	}
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
