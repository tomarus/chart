package chart

import "time"

// xlabels creates n time labels for use on the x axis.
func (c *Chart) xlabels(n int) []string {
	skip := (c.end - c.start) / int64(n)
	tdsave := ""
	res := []string{}
	b := c.start
	for ix := 0; ix <= c.width; ix += c.width / n {
		if ix == 0 || ix == c.width {
			b += skip
			continue
		}
		td := time.Unix(b, 0).Format("01-02")
		txt := time.Unix(b, 0).Format("15:04")
		if td != tdsave {
			tdsave = td
			txt = time.Unix(b, 0).Format("01-02 15:04")
		}
		res = append(res, txt)
		b += skip
	}
	return res
}

// scales creates all scale labels for all datasets.
func (c *Chart) scales(n int) {
	for i := range c.data {
		c.data[i].CreateScale(n)
	}
}
