package axis

import (
	"time"

	"github.com/tomarus/chart/format"
	"github.com/tomarus/chart/image"
)

// Axis defines an axis (doh)
type Axis struct {
	position Position
	format   Formatter
	duration time.Duration
	grid     int
	ticks    int
	center   bool
}

// Formatter is the callback interface function used to format a label.
type Formatter func(value float64) string

// Liner is the callback interface function to draw a line.
// Any padding/margin offsets should be applied by the caller.
type Liner func(x1, y1, x2, y2 int)

// Position defines the position for the axis. Bottom, Top, Left or Right
type Position int

const (
	// Bottom defines an axis aligned to the bottom of the image.
	Bottom Position = iota

	// Left defines an axis aligned to the let of the image.
	// Only Bottom and Left are defined at the moment.
	Left
)

// New creates a new Axis on the specified position using the formatter.
func New(p Position, f Formatter) *Axis {
	return &Axis{
		position: p,
		format:   f,
	}
}

// NewSI creates a new Axis on the specified position using the default SI formatter.
func NewSI(p Position, base int) *Axis {
	return New(p, func(in float64) string {
		return format.SI(in, 1, float64(base), "", "", "")
	})
}

// NewTime creates a new Axis on the specified position using the default Time formatter.
// A timefmt is specified using the default Go Time format, e.g. 2006-01-02 15:04
func NewTime(p Position, timefmt string) *Axis {
	return New(p, func(in float64) string {
		return time.Unix(int64(in), 0).Format(timefmt)
	})
}

// Ticks sets the number of gridlines/labels or ticks for this axis.
// This configures an equally centered grid pattern for an axis.
// Use either one of Ticks() or Duration()
func (a *Axis) Ticks(n int) *Axis {
	a.ticks = n
	return a
}

// Duration sets the time period for the gridlines/labels or ticks for the axis.
// This configures a grid aligned to the nearest time unit.
// Use either one of Ticks() or Duration()
func (a *Axis) Duration(d time.Duration) *Axis {
	a.duration = d
	return a
}

// Grid displays a grid for this axis. By default it doesn't show a grid.
// n is the amount of gridlines to draw between axis ticks. If your axis
// spans 4 hours and n = 2 then a gridline will be drawn every 2 hours.
func (a *Axis) Grid(n int) *Axis {
	a.grid = n
	return a
}

// Center aligns the label in tne center of the grid instead of the start/end of grid.
func (a *Axis) Center() *Axis {
	a.center = true
	return a
}

// Draw renders the grid and labels.
// mx/my is the top-left start position, the margin (or offset).
// FIXME there are text-margin constants in this function which are probably dependent on the font and size used.
// It also depends a bit too much on the actual font/line/color drawing stuff.
func (a *Axis) Draw(img image.Image, w, h, mx, my int, min, max float64) {
	const col = "title2"

	switch a.position {
	case Bottom:
		off := 0 // grid line offset, TODO calculate from rounded time.Duration
		if a.duration > 0 {
			t := int(a.duration / time.Second)
			a.ticks = 1
			if t > 0 {
				a.ticks = int(max-min) / t
			}
		}
		if a.grid > 0 {
			t := a.ticks * a.grid
			o := w / t
			for dx := o; dx < w; dx += o {
				col := "grid"
				if dx/o%a.grid != 0 {
					col = "grid2"
				}
				img.Line(col, dx+mx+off, my, dx+mx+off, my+h)
			}
		}

		toff := 0 // text offset
		if a.center {
			toff = w / a.ticks / 2
		}
		for dx := (w / a.ticks) - off; dx < w+toff; dx += w / a.ticks {
			str := a.format(min + ((max-min)/float64(w))*(float64(dx)-float64(toff)))
			img.TextID("grid", col, "middle", image.GridRole, dx+mx+off-toff, my+h+14, str) // FIXME "14" (padding/offset)
		}
	case Left:
		if a.grid > 0 {
			t := a.ticks * a.grid
			o := h / t
			for dy := o; dy < h; dy += o {
				col := "grid"
				if dy/o%a.grid != 0 {
					col = "grid2"
				}
				img.Line(col, mx, dy+my, mx+w, dy+my)
			}
		}
		sc := a.Scales(h, min, max)
		i := 0
		// TODO if show or hide zero value option
		for dy := 0; dy <= h; dy += h / a.ticks {
			str := sc[i]
			i++
			img.TextID("ygrid", col, "end", image.GridRole, mx-4, dy+my+4, str) // FIXME "4" (padding/spacing)
		}
	}
}

// Scales creates and formats the Y-axis scale.
func (a *Axis) Scales(h int, min, max float64) []string {
	s := []string{}
	switch a.position {
	case Left:
		for dy := 0; dy <= h; dy += h / a.ticks {
			str := a.format(max - ((max-min)/float64(h))*float64(dy))
			s = append(s, str)
		}
	}
	return s
}
