// Package svg provides the svg interface for tomarus chart lib.
package svg

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tomarus/chart/data"
	"github.com/tomarus/chart/palette"
)

type SVG struct {
	w                io.Writer
	width, height    int
	marginx, marginy int
	start, end       int
	pal              *palette.Palette
}

func New() *SVG {
	return &SVG{}
}

func (svg *SVG) Start(wr io.Writer, w, h, mx, my, start, end int, p *palette.Palette) {
	svg.w = wr
	svg.width = w
	svg.height = h
	svg.marginx = mx
	svg.marginy = my
	svg.start = start
	svg.end = end
	svg.pal = p

	svg.svgHead(w+mx+4, h+(2*my))
	svg.svgCSS(svg.pal)
	svg.p(`<rect class="background" x="0" y="0" width="%d" height="%d"/>`, w+mx+32, h+(2*my))
}

func (svg *SVG) Graph(d data.Collection) {
	svg.p(`<defs>`)

	svg.p(`<script type="text/javascript"><![CDATA[`)
	svg.p("const w=%d,h=%d,mx=%d,my=%d,start=%d,end=%d", svg.width, svg.height, svg.marginx, svg.marginy, svg.start, svg.end)
	jsdata, _ := json.Marshal(d)
	svg.p("const data=" + string(jsdata))
	fmt.Fprint(svg.w, js)
	svg.p("]]></script>")

	for i := range d {
		svg.p(`<g id="path%d">`, i+1)
		svg.p(`<path style="fill: none; stroke: %s; shape-rendering: crispEdges" d="M0,0"/>`, svg.pal.GetHexAxisColor(i))
		svg.p(`</g>`)
	}

	fmt.Fprintln(svg.w, `</defs>`)

	for i := range d {
		svg.p(`<use x="%d" y="%d" xlink:href="#path%d"/>`, svg.marginx, svg.marginy, i+1)
	}

	svg.p(`<line id="markerx" x1="0" x2="0" y1="%d" y2="%d" class="marker" style="visibility:hidden"/>`, svg.marginy, svg.height+svg.marginy)
	svg.p(`<line id="markery" x1="%d" x2="%d" y1="0" y2="0" class="marker" style="visibility:hidden"/>`, svg.marginx, svg.width+svg.marginx)
	svg.p(`<line id="markerx2" x1="0" x2="0" y1="%d" y2="%d" class="marker" style="visibility:hidden"/>`, svg.marginy, svg.height+svg.marginy)
	svg.p(`<line id="markery2" x1="%d" x2="%d" y1="0" y2="0" class="marker" style="visibility:hidden"/>`, svg.marginx, svg.width+svg.marginx)
	svg.p(`<rect id="markersel" x="0" y="0" width="0" height="0" class="" style='fill-opacity:.5;fill:%s'/>`, svg.pal.GetHexColor("select"))
	svg.p(`<text class="title" id="markertext" x="%d" y="%d" />`, svg.marginx, svg.marginy/2+4)

	svg.drawMA()
}

func (svg *SVG) drawMA() {
	const maColor = "color3"
	fmt.Fprintln(svg.w, `<defs>`)
	svg.p(`<g id="ma">`)
	svg.p(`<path style="fill: none; stroke: %s; stroke-width: 2; shape-rendering: auto" d="M0,0"/>`, svg.pal.GetHexColor(maColor))
	svg.p(`</g>`)
	fmt.Fprintln(svg.w, `</defs>`)
	svg.p(`<use x="%d" y="%d" xlink:href="#ma"/>`, svg.marginx, svg.marginy)

	y := svg.height + svg.marginy + 4
	svg.p(`<rect id="mabut" x="%d" y="%d" width="12" height="12" style="visibility:normal;fill:%s"/>`, svg.width+svg.marginx-12, y, svg.pal.GetHexColor(maColor))
}

func (svg *SVG) Text(color, align string, x, y int, txt string) {
	anchor := ""
	switch align {
	case "begin", "left":
		anchor = "text-anchor:left"
	case "middle", "center":
		anchor = "text-anchor:middle"
	case "end", "right":
		anchor = "text-anchor:end"
	}
	svg.p(`<g class="%s"><text style="%s" x="%d" y="%d">%s</text></g>`, color, anchor, x, y, txt)
}

func (svg *SVG) ID(id string) {
	svg.p(`<g id="%s">`, id)
}

func (svg *SVG) EndID() {
	svg.p(`</g>`)
}

func (svg *SVG) Line(color string, x1, y1, x2, y2 int) {
	svg.p(`<line class="%s" x1="%d" y1="%d" x2="%d" y2="%d"/>`, color, x1, y1, x2, y2)
}

// Legend writes the legend buttons to the bottom.
func (svg *SVG) Legend(d data.Collection, p *palette.Palette) {
	x := svg.marginx
	y := svg.height + svg.marginy + 4
	svg.p(`<g class="legend">`)
	for i := range d {
		id := fmt.Sprintf("path%d_b", i+1)
		svg.p(`<g id="%s"><rect x="%d" y="%d" width="12" height="12" style="visibility:normal;fill:%s"/></g>`, id, x, y, p.GetHexAxisColor(i))
		x += 16
	}
	svg.p(`</g>`)
}

func (svg *SVG) End() error {
	svg.p(`</svg>`)
	return nil
}

func (svg *SVG) Border(x, y, w, h int) {
	svg.svgRect(x, y, w, h, "border")
}

func (svg *SVG) svgRect(x, y, w, h int, class string) {
	svg.p(`<rect x="%d" y="%d" width="%d" height="%d" class="%s"/>`, x, y, w, h, class)
}

func (svg *SVG) svgHead(w, h int) {
	svg.p(`<?xml version="1.0"?>`)
	svg.p(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">"`, w, h, w, h)
}

func (svg *SVG) svgCSS(p *palette.Palette) {
	svg.p(`<defs><style type="text/css"><![CDATA[`)
	svg.p(".grid { stroke: %s; stroke-opacity: .62; stroke-dasharray: 1; stroke-width: 0.25 }", p.GetHexColor("grid"))
	svg.p(".title { fill: %s; font-size: 11px; font-family: menlo; fill-opacity:1; stroke-width: 0 }", p.GetHexColor("title"))
	svg.p(".border { stroke: %s; stroke-opacity: 1; stroke-width: 1; fill: none }", p.GetHexColor("border"))
	svg.p(".marker { stroke: %s; stroke-opacity: 1; stroke-width: 1 }", p.GetHexColor("marker"))
	svg.p(".background { fill: %s }", p.GetHexColor("background"))
	svg.p(".legend { cursor: pointer }")
	svg.p("]]></style></defs>")
}

func (svg *SVG) p(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(svg.w, format+"\n", a...)
}
