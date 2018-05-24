// Package svg provides the svg interface for tomarus chart lib.
package svg

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tomarus/chart/data"
	"github.com/tomarus/chart/format"
	"github.com/tomarus/chart/image"
	"github.com/tomarus/chart/palette"
)

// SVG implements the chart interface to write SVG images.
type SVG struct {
	w                io.Writer
	data             data.Collection
	width, height    int
	marginx, marginy int
	start, end       int64
	pal              *palette.Palette
	txtids           map[string][]textid
}

type textid struct {
	color, align, txt string
	x, y              int
	role              image.TextRole
}

// New initializes a new svg chart image writer.
func New() *SVG {
	return &SVG{txtids: make(map[string][]textid)}
}

// Start initializes a new image and sets the defaults.
func (svg *SVG) Start(wr io.Writer, w, h, mx, my int, start, end int64, p *palette.Palette, d data.Collection) {
	svg.w = wr
	svg.data = d
	svg.width = w
	svg.height = h
	svg.marginx = mx
	svg.marginy = my
	svg.start = start
	svg.end = end
	svg.pal = p

	svg.svgHead(w+mx+4, h+(2*my)+((d.Len()+1)*16))
	svg.svgCSS(svg.pal)
	svg.p(`<rect class="background" x="0" y="0" width="%d" height="%d"/>`, w+mx+32, h+(2*my)+((d.Len()+1)*16))
}

// Graph renders all chart dataset values to the visible chart area.
func (svg *SVG) Graph() error {
	svg.p(`<defs>`)
	{
		svg.p(`<script type="text/javascript"><![CDATA[`)
		svg.p("const w=%d,h=%d,mx=%d,my=%d,start=%d,end=%d", svg.width, svg.height, svg.marginx, svg.marginy, svg.start*1000, svg.end*1000)
		jsdata, _ := json.Marshal(svg.data)
		svg.p("const data=" + string(jsdata))
		fmt.Fprint(svg.w, js)
		svg.p("]]></script>")

		for i := range svg.data {
			svg.p(`<g id="path%d">`, i+1)
			svg.p(`<path style="fill: none; stroke: %s; shape-rendering: crispEdges" d="M0,0"/>`, svg.pal.GetHexAxisColor(i))
			svg.p(`</g>`)
		}
	}
	fmt.Fprintln(svg.w, `</defs>`)

	for i := range svg.data {
		svg.p(`<use x="%d" y="%d" xlink:href="#path%d"/>`, svg.marginx, svg.marginy, i+1)
	}

	svg.p(`<line id="markerx" x1="0" x2="0" y1="%d" y2="%d" class="marker" style="visibility:hidden"/>`, svg.marginy, svg.height+svg.marginy)
	svg.p(`<line id="markery" x1="%d" x2="%d" y1="0" y2="0" class="marker" style="visibility:hidden"/>`, svg.marginx, svg.width+svg.marginx)
	svg.p(`<line id="markerx2" x1="0" x2="0" y1="%d" y2="%d" class="marker" style="visibility:hidden"/>`, svg.marginy, svg.height+svg.marginy)
	svg.p(`<line id="markery2" x1="%d" x2="%d" y1="0" y2="0" class="marker" style="visibility:hidden"/>`, svg.marginx, svg.width+svg.marginx)
	svg.p(`<rect id="markersel" x="0" y="0" width="0" height="0" class="" style='fill-opacity:.25;fill:%s'/>`, svg.pal.GetHexColor("select"))
	svg.p(`<g class="title gridfont"><text id="markertext" x="%d" y="%d"/></g>`, svg.marginx, svg.marginy/2+4)

	svg.drawMA()
	return nil
}

func (svg *SVG) drawMA() {
	const maColor = "marker"
	svg.p(`<defs>`)
	svg.p(`<g id="ma">`)
	svg.p(`<path style="fill: none; stroke: %s; stroke-width: 2; shape-rendering: auto" d="M0,0"/>`, svg.pal.GetHexColor(maColor))
	svg.p(`</g>`)
	svg.p(`</defs>`)
	svg.p(`<use x="%d" y="%d" xlink:href="#ma"/>`, svg.marginx, svg.marginy)

	y := svg.height + svg.marginy + 4
	svg.p(`<rect id="mabut" x="%d" y="%d" width="12" height="12" style="visibility:normal;fill:%s"/>`, svg.width+svg.marginx-12, y, svg.pal.GetHexColor(maColor))
}

// Text writes a string to the image.
func (svg *SVG) Text(color, align string, role image.TextRole, x, y int, txt string) {
	anchor := ""
	switch align {
	case "begin", "left":
		anchor = "text-anchor:left"
	case "middle", "center":
		anchor = "text-anchor:middle"
	case "end", "right":
		anchor = "text-anchor:end"
	}
	class := color
	if role == image.TitleRole {
		class += " titlefont"
	} else {
		class += " gridfont"

	}
	svg.p(`<g class="%s"><text style="%s" x="%d" y="%d">%s</text></g>`, class, anchor, x, y, txt)
}

// TextID writes a string to the image using an id.
func (svg *SVG) TextID(id, color, align string, role image.TextRole, x, y int, txt string) {
	if svg.txtids[id] == nil {
		svg.txtids[id] = make([]textid, 0)
	}
	svg.txtids[id] = append(svg.txtids[id], textid{color, align, txt, x, y, role})
}

func (svg *SVG) drawTextIDs() {
	for k, v := range svg.txtids {
		svg.p(`<g id="%s">`, k)
		for _, t := range v {
			svg.Text(t.color, t.align, t.role, t.x, t.y, t.txt)
		}
		svg.p(`</g>`)
	}
}

// Line draws a line between the points using the color name from the palette.
func (svg *SVG) Line(color string, x1, y1, x2, y2 int) {
	svg.p(`<line class="%s" x1="%d" y1="%d" x2="%d" y2="%d"/>`, color, x1, y1, x2, y2)
}

// Legend writes the legend buttons to the bottom.
func (svg *SVG) Legend(base float64) {
	x := svg.marginx
	y := svg.height + svg.marginy + 4 + 16

	q := "Min     Max     Avg"
	svg.Text("title2", "right", image.GridRole, x+svg.width, y+11, q)
	y += 16

	for i, d := range svg.data {
		id := fmt.Sprintf("path%d_b", i+1)
		svg.p(`<g id="%s" class="legend"><rect x="%d" y="%d" width="12" height="12" style="visibility:normal;fill:%s"/></g>`, id, x, y, svg.pal.GetHexAxisColor(i))
		svg.Text("title", "left", image.GridRole, x+20, y+11, d.Title)

		// FIXME use axis formatters for this.
		min, max, avg := d.MinMaxAvg()
		mmax := format.SI(max, 1, base, "", "", "")
		mmin := format.SI(min, 1, base, "", "", "")
		mavg := format.SI(avg, 1, base, "", "", "")
		q := fmt.Sprintf("%6s  %6s  %6s", mmin, mmax, mavg)
		svg.Text("title", "right", image.GridRole, x+svg.width, y+11, q)
		svg.Line("grid2", x, y+11+3, x+svg.width, y+11+3)
		y += 16
	}
}

// End finishes and writes the image to the output writer.
func (svg *SVG) End() error {
	svg.drawTextIDs()
	svg.p(`</svg>`)
	return nil
}

// Border draws a border around the chart area.
func (svg *SVG) Border(x, y, w, h int) {
	svg.svgRect(x, y, w, h, "border")
}

func (svg *SVG) svgRect(x, y, w, h int, class string) {
	svg.p(`<rect x="%d" y="%d" width="%d" height="%d" class="%s"/>`, x, y, w, h, class)
}

func (svg *SVG) svgHead(w, h int) {
	svg.p(`<?xml version="1.0"?>`)
	svg.p(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">`, w, h, w, h)
}

func (svg *SVG) svgCSS(p *palette.Palette) {
	svg.p(`<defs><style type="text/css"><![CDATA[`)
	svg.p("* { shape-rendering: crispEdges; }")
	svg.p(".grid { stroke: %s; stroke-opacity: %f; stroke-dasharray: 1; stroke-width: .75 }", p.GetHexColor("grid"), p.GetAlpha("grid"))
	svg.p(".grid2 { stroke: %s; stroke-opacity: %f; stroke-dasharray: 1; stroke-width: .33 }", p.GetHexColor("grid2"), p.GetAlpha("grid2"))

	svg.p(".title { fill: %s; fill-opacity: .75 }", p.GetHexColor("title"))
	svg.p(".title2 { fill: %s; fill-opacity: .75 }", p.GetHexColor("title2"))
	svg.p(".titlefont { font-variant: small-caps; font-style: italic; font-size: 18px; font-family: menlo; }")
	svg.p(".gridfont { font-size: 13px; font-family: menlo; stroke-width: .33; }")

	svg.p(".border { stroke: %s; stroke-opacity: .666; fill: none }", p.GetHexColor("border"))
	svg.p(".marker { stroke: %s; stroke-opacity: 1; stroke-width: 1; fill: none; }", p.GetHexColor("marker"))
	svg.p(".background { fill: %s }", p.GetHexColor("background"))
	svg.p(".legend { cursor: pointer }")
	svg.p("text { white-space: pre }")
	svg.p("]]></style></defs>")
}

func (svg *SVG) p(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(svg.w, format+"\n", a...)
}
