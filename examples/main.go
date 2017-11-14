package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/tomarus/chart"
	"github.com/tomarus/chart/axis"
	"github.com/tomarus/chart/data"
	"github.com/tomarus/chart/image"
	"github.com/tomarus/chart/png"
	"github.com/tomarus/chart/svg"
)

func main() {
	listen := flag.String("listen", ":3000", "Address for HTTP listener")
	flag.Parse()

	http.HandleFunc("/chart.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		drawChart(w, r, svg.New(), iFormValue(r, "width"), fFormValue(r, "h1"), fFormValue(r, "h2"), fFormValue(r, "add1"), fFormValue(r, "add2"))
	})
	http.HandleFunc("/chart.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		drawChart(w, r, png.New(), iFormValue(r, "width"), fFormValue(r, "h1"), fFormValue(r, "h2"), fFormValue(r, "add1"), fFormValue(r, "add2"))
	})

	http.HandleFunc("/chartsmall.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		drawChartSmall(w, r, svg.New())
	})
	http.HandleFunc("/chartsmall.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		drawChartSmall(w, r, png.New())
	})

	http.HandleFunc("/chartthemed.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		drawChartThemed(w, r, svg.New(), r.FormValue("theme"), r.FormValue("scheme"))
	})
	http.HandleFunc("/chartthemed.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		drawChartThemed(w, r, png.New(), r.FormValue("theme"), r.FormValue("scheme"))
	})

	http.HandleFunc("/", html)

	log.Printf("Listening on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func fFormValue(req *http.Request, name string) float64 {
	x := req.FormValue(name)
	r, _ := strconv.ParseFloat(x, 64)
	return r
}

func iFormValue(req *http.Request, name string) int64 {
	x := req.FormValue(name)
	r, _ := strconv.ParseInt(x, 10, 64)
	return r
}

func html(w http.ResponseWriter, r *http.Request) {
	p := r.FormValue("p")
	fmt.Fprint(w, htmltop)
	switch p {
	case "basic":
		fmt.Fprint(w, htmlbasic)
	case "colors":
		fmt.Fprint(w, htmlcolors)
	case "random":
		fmt.Fprint(w, htmlrandom)
	}
	fmt.Fprint(w, htmlbot)
}

const htmltop = `
<html>
<body>
<h1>Chart Examples</h1>
<ul>
	<li><a href="?p=basic">Basic Charts</a></li>
	<li><a href="?p=colors">Builtin Themes</a></li>
	<li><a href="?p=random">Random Themes</a></li>
</ul>
`

const htmlbasic = `
<div>
	<h2>Basic Charts</h2>
	<object data="/chart.svg?h1=256&h2=128&add1=1&add2=1"></object>
	<img src="/chart.png?h1=256&h2=128&add1=1&add2=1"></img>
</div>
<div>
	<h2>Small Value Charts</h2>
	<object data="/chart.svg?h1=1&h2=0.5&add1=3&add2=3"></object>
	<img src="/chart.png?h1=1&h2=0.5&add1=3&add2=3"></img>
</div>
<div>
	<h2>Few Values Charts</h2>
	<object data="/chartsmall.svg"></object>
	<img src="/chartsmall.png"></img>
</div>
`

const htmlcolors = `
<div>
	<h2>Themed Charts</h2>
	<h3>Scheme "white"</h3>
	<object data="/chartthemed.svg?scheme=white"></object>
	<img src="/chartthemed.png?scheme=white"></img>
	<h3>Scheme "black"</h3>
	<object data="/chartthemed.svg?scheme=black"></object>
	<img src="/chartthemed.png?scheme=black"></img>
	<h3>Scheme "pink"</h3>
	<object data="/chartthemed.svg?scheme=pink"></object>
	<img src="/chartthemed.png?scheme=pink"></img>
	<h3>Scheme "solarized"</h3>
	<object data="/chartthemed.svg?scheme=solarized"></object>
	<img src="/chartthemed.png?scheme=solarized"></img>
</div>
`

const htmlrandom = `
<div>
	<h3>Scheme "random" theme "dark"</h3>
	<object data="/chartthemed.svg?theme=dark&scheme=random"></object>
	<img src="/chartthemed.png?theme=dark&scheme=random"></img>
	<h3>Scheme "random" theme "light"</h3>
	<object data="/chartthemed.svg?theme=light&scheme=random"></object>
	<img src="/chartthemed.png?theme=light&scheme=random"></img>

	<h3>Scheme "hsl:0,0.66,0.5" theme "light"</h3>
	<object data="/chartthemed.svg?theme=light&scheme=hsl:0,0.66,0.5"></object>
	<img src="/chartthemed.png?theme=light&scheme=hsl:0,0.66,0.5"></img>
	<h3>Scheme "hsl:320,0.25,0.5" theme "light"</h3>
	<object data="/chartthemed.svg?theme=light&scheme=hsl:320,0.25,0.5"></object>
	<img src="/chartthemed.png?theme=light&scheme=hsl:320,0.25,0.5"></img>
	<h3>Scheme "hsl:120,0.5,0.4" theme "dark"</h3>
	<object data="/chartthemed.svg?theme=dark&scheme=hsl:120,0.5,0.4"></object>
	<img src="/chartthemed.png?theme=dark&scheme=hsl:120,0.5,0.4"></img>
</div>
`

const htmlbot = `
</body>
</html>
`

func mksin(width int64, height float64, freq int64, add float64) []float64 {
	res := width / freq
	fx := 0.
	sin := make([]float64, width)
	for i := int64(0); i < width; i++ {
		ent := (math.Sin(fx) + add) * float64(height)
		sin[i] = ent
		fx += (2 * math.Pi) / float64(res)
	}
	return sin
}

func drawChart(w http.ResponseWriter, r *http.Request, img image.Image, width int64, h1, h2, add1, add2 float64) {
	if width == 0 {
		width = 720
	}
	opts := &chart.Options{
		Title:  "Example Chart",
		Image:  img,
		Size:   "small",
		Scheme: "random",
		Theme:  "light",
		Start:  time.Now().AddDate(0, 0, -2).Unix(),
		End:    time.Now().Unix(),
		W:      w,
		Axes: []*axis.Axis{
			axis.NewTime(axis.Bottom, "Mon 15:04").Duration(8 * time.Hour).Grid(4),
			axis.NewSI(axis.Left, 1000).Ticks(4).Grid(2),
		},
		// Data: []data.Data{
		// 	data.New(mksin(width, h1, 2, add1)).Title("dataset 1").Type("area").Gap(0),
		// 	data.New(mksin(width, h2, 4, add2)).Title("dataset 2").Type("line").Gap(0),
		// },
	}

	c, err := chart.NewChart(opts)
	if err != nil {
		panic(err)
	}

	c.AddData(&data.Options{Type: "area", Title: "dataset 1", Gap: 0}, mksin(width, h1, 2, add1))
	c.AddData(&data.Options{Type: "line", Title: "dataset 2", Gap: 0}, mksin(width, h2, 4, add2))
	c.Render()
}

func drawChartSmall(w http.ResponseWriter, r *http.Request, img image.Image) {
	opts := &chart.Options{
		Title:  "Example Chart",
		Image:  img,
		Size:   "small",
		Scheme: "random",
		Theme:  "light",
		Start:  time.Now().AddDate(0, 0, -30).Unix(),
		End:    time.Now().Unix(),
		W:      w,
		Axes: []*axis.Axis{
			axis.NewTime(axis.Bottom, "02").Duration(1 * 86400 * time.Second).Grid(1).Center(),
			axis.NewSI(axis.Left, 1000).Ticks(10).Grid(1),
		},
	}

	c, _ := chart.NewChart(opts)
	c.AddData(&data.Options{Type: "area", Title: "Dataset 1", Gap: .05}, mksin(30, 256, 1, 3))
	c.AddData(&data.Options{Type: "area", Title: "Dataset Nr 2", Gap: .05}, mksin(30, 128, 2, 1))
	c.Render()
}

func drawChartThemed(w http.ResponseWriter, r *http.Request, img image.Image, theme, scheme string) {
	opts := &chart.Options{
		Title:  "Example Chart",
		Image:  img,
		Size:   "small",
		Scheme: scheme,
		Theme:  theme,
		Start:  time.Now().AddDate(0, 0, -2).Unix(),
		End:    time.Now().Unix(),
		W:      w,
		Axes: []*axis.Axis{
			axis.NewTime(axis.Bottom, "Mon 15:04").Duration(8 * time.Hour).Grid(4),
			axis.NewSI(axis.Left, 1000).Ticks(4).Grid(2),
		},
		// Data: []data.Data{
		// 	data.New(mksin(width, h1, 2, add1)).Title("dataset 1").Type("area").Gap(0),
		// 	data.New(mksin(width, h2, 4, add2)).Title("dataset 2").Type("line").Gap(0),
		// },
	}

	c, err := chart.NewChart(opts)
	if err != nil {
		panic(err)
	}

	c.AddData(&data.Options{Title: "dataset 1", Gap: 0}, mksin(720, 4, 2, 1))
	c.AddData(&data.Options{Title: "dataset 2", Gap: 0}, mksin(720, 3, 7, 1))
	c.AddData(&data.Options{Title: "dataset 3", Gap: 0}, mksin(720, 2, 3, 1))
	c.AddData(&data.Options{Title: "dataset 4", Gap: 0}, mksin(720, 1, 11, 1))
	c.Render()
}
