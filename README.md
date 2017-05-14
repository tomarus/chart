# Go Chart Lib

Dead simple rrd like interactive svg graphs with focus on pixel perfect rendering of source data.

Written in Go, the output is a SVG image containing Javascript which does most rendering and selection magic.

## Examples

Example screenshot:

[View as interactive SVG](http://s.chiparus.org/5/5989676a301be238.svg)

![Example Screenshot](http://s.chiparus.org/5/5860a66293f1a6f1.png)

Example of useless but interesting random colors:

![Example of totally useless random colors](http://s.chiparus.org/7/7b2fd43470e2475b.png)

## Example Usage

```go
import "github.com/tomarus/chart"
opts := &chart.Options{
    Title:  "Traffic",
    Size:   "big",   // big is 1440px, small is 720px
    Scheme: "white", // or black/random/pink/solarized or hsv:180,0.5,0.25
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

It currently has no Go and no Javascript dependencies. But they will be added eventually probably.

This project has just started and a lot of stuf is still missing or incomplete.

This is a small list of ideas/todos:
* Data interpolation (data must match graph size now)
* Flexible graph sizes
* Bar charts (basically wider area pixels)
* Custom lines and markers, like 95th percentile line, downtime markers, etc
* Finish PNG implementation
