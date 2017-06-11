package chart

// xgrid lays out the grid on the x axis and creates the labels.
func (c *Chart) xgrid(x, y, n int, bt, et int64) {
	off := (c.width / n) + x
	for ix := 0; ix < n-1; ix++ {
		nx := off + ((c.width / n) * ix)
		c.image.Line("grid", nx, y, nx, y+c.height)
	}

	labels := c.xlabels(n) // XXX
	for ix := 0; ix < n-1; ix++ {
		nx := off + ((c.width / n) * ix)
		c.image.Text("title", "middle", nx, y+c.height+14, labels[ix])
	}
}

// ygrid lays out the grid on the y axis and creates the label attributes.
func (c *Chart) ygrid(x, y, n int) {
	off := (c.height / n) + y
	for iy := 0; iy < n-1; iy++ {
		ny := off + ((c.height / n) * iy)
		c.image.Line("grid", x, ny, x+c.width, ny)
	}

	c.image.ID("ygrid")
	scale := 0
	for iy := y; iy <= c.height; iy += c.height / n {
		c.image.Text("title", "end", x-4, iy+4, c.data[0].Scale[scale])
		scale++
	}
	c.image.EndID()
}
