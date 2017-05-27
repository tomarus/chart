# Go Chart Lib

Dead simple rrd like interactive svg graphs with focus on pixel perfect rendering of source data.

Written in Go, the output is a SVG image containing Javascript which does most rendering and selection magic.
Additionally a static png image can be written.

## Examples

Example screenshot:

[View as interactive SVG](http://s.chiparus.org/6/6b15c5349e894fe9.svg)

![Example Spngcreenshot](http://s.chiparus.org/5/5caa4e08e4b2edb3.png)

Example of useless but interesting random colors:

![Example of totally useless random colors](http://s.chiparus.org/7/7b2fd43470e2475b.png)

## Example Usage

```go
import (
    "github.com/tomarus/chart"
    "github.com/tomarus/chart/svg"
)
opts := &chart.Options{
    Title:  "Traffic",
    Image:  svg.New(), // or png.New()
    Size:   "auto",    // big is 1440px, small is 720px, auto is size of dataset
    Height: 300,       // Defaults to -1, when size=auto height=width/4, otherwise set fixed height
    Scheme: "white",   // or black/random/pink/solarized or hsl:180,0.5,0.25
    Start:  start_epoch,
    End:    end_epoch,
    Xdiv:   12,
    Ydiv:   5,
    W:      w,
}
c, err := chart.NewChart(opts)
if err != nil {
    panic(err)
}
warn := c.AddData("area", []yourData)
if err != nil {
    fmt.Println(warn)
}
w.Header().Set("Content-Type", "image/svg+xml")
c.Render()
```

## Notes

This is an experimental work in progress for educational and research purposes only.

This project has very few Go dependencies and no Javascript dependencies. Only freetype and image/font is used for png output.

This project has just started and a lot of stuf is still missing or incomplete.

This is a small list of ideas/todos:
* Custom lines and markers, like 95th percentile line, downtime markers, etc
* Complete legends and dataset descriptions