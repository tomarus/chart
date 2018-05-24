package image

import (
	"io"

	"github.com/tomarus/chart/data"
	"github.com/tomarus/chart/palette"
)

// TextRole is used to distinguish between fonts.
type TextRole int

const (
	// TitleRole is used to identify title text
	TitleRole TextRole = iota

	// GridRole is used to identify grid or scale text
	GridRole
)

// Image defines the interface for image (svg/png) backends.
type Image interface {
	// Start initializes a new image and sets the defaults.
	Start(wr io.Writer, w, h, mx, my int, start, end int64, p *palette.Palette, d data.Collection)

	// End finishes and writes the image to the output writer.
	End() error

	// Graph renders all chart dataset values to the visible chart area.
	Graph() error

	// Text writes a string to the image.
	Text(color, align string, role TextRole, x, y int, txt string)

	// TextID writes a string to the image using an id.
	TextID(id, color, align string, role TextRole, x, y int, txt string)

	// Line draws a line between two points. The current implementations need to
	// only draw horizontal or vertical lines.
	Line(color string, x1, y1, x2, y2 int)

	// Legend draws the legend in the image using base for SI formatting.
	Legend(float64)

	// Border draws a border around the chart area.
	Border(x, y, w, h int)
}
