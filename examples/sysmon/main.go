// Command sysmon is a basic /proc system monitoring tool storing only
// the last 15 minutes of data, sampled each second, in memory.
// It monitors and plots usage data for cpu, free mem,
// network, load avg, procs and disk io.
// It's an example using the tomarus/chart and c9s/goprocinfo packages.
package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/tomarus/chart"
	"github.com/tomarus/chart/axis"
	"github.com/tomarus/chart/data"
	"github.com/tomarus/chart/examples/sysmon/mods"
	"github.com/tomarus/chart/png"
)

const delay = 1 * time.Second

var inputs = []mods.Collector{
	&mods.CPUStat{},
	&mods.MemInfo{},
	&mods.NetDev{},
	&mods.LoadAvg{},
	&mods.Procs{},
	&mods.DiskStat{},
}

func collector() {
	t := time.NewTicker(delay)
	for range t.C {
		for i := range inputs {
			go func(i int) {
				err := inputs[i].Update()
				if err != nil {
					log.Printf("Update %v: %v", inputs[i], err)
				}
			}(i)
		}
	}
}

func plot(input mods.Collector, title string, w http.ResponseWriter, r *http.Request) {
	L := input.Len()
	dur := (time.Duration(L) * time.Second) / 5
	opts := &chart.Options{
		Title:  title,
		Image:  png.New(),
		Width:  900,
		Height: 300,
		Scheme: fmt.Sprintf("hsl:%d,0.5,0.5", hashN(title, 120)),
		Theme:  "light",
		Start:  time.Now().Add(-time.Duration(L) * time.Second).Unix(),
		End:    time.Now().Unix(),
		W:      w,
		Axes: []*axis.Axis{
			axis.NewTime(axis.Bottom, "15:04:05").Duration(dur).Grid(3),
			axis.NewSI(axis.Left, 1000).Ticks(4).Grid(2),
		},
	}

	ch, err := chart.NewChart(opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("NewChart Error: %v", err)
		return
	}

	for _, v := range input.Data() {
		err := ch.AddData(&data.Options{Title: v.Title}, v.Values)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Data Error: %v", err)
			return
		}
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Content-Type", "image/png")
	err = ch.Render()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Render Error: %v", err)
	}
}

func hashN(s string, n uint32) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() % n
}

var html = `
<html>
<head>
<script type="text/javascript">
window.onload = function() {
    function updateImage(id) {
        var img = document.getElementById(id);
        img.src = img.src.split("?")[0] + "?" + new Date().getTime();
    }
    var charts = ["mem", "cpu", "net", "load", "proc", "io"]
    const intvl = ` + fmt.Sprintf("%d", delay/time.Millisecond) + `
    charts.map((chart, i)=>{
        setTimeout(()=>{
            setInterval(()=>{
                    updateImage(chart)
            }, intvl)
        }, i*(intvl/charts.length))
    })
}
</script>
</head>
<body>
<img align="top" id="load" src="/load.png"></img>
<img align="top" id="cpu" src="/cpu.png"></img>
<br/><p/><br/><p/>
<img align="top" id="mem" src="/mem.png"></img>
<img align="top" id="net" src="/net.png"></img>
<br/><p/><br/><p/>
<img align="top" id="proc" src="/proc.png"></img>
<img align="top" id="io" src="/io.png"></img>
</body>
</html>
`

func main() {
	go collector()

	http.HandleFunc("/cpu.png", func(w http.ResponseWriter, r *http.Request) {
		plot(inputs[0], "CPU", w, r)
	})
	http.HandleFunc("/mem.png", func(w http.ResponseWriter, r *http.Request) {
		plot(inputs[1], "MEM", w, r)
	})
	http.HandleFunc("/net.png", func(w http.ResponseWriter, r *http.Request) {
		plot(inputs[2], "NET", w, r)
	})
	http.HandleFunc("/load.png", func(w http.ResponseWriter, r *http.Request) {
		plot(inputs[3], "LOAD", w, r)
	})
	http.HandleFunc("/proc.png", func(w http.ResponseWriter, r *http.Request) {
		plot(inputs[4], "PROC", w, r)
	})
	http.HandleFunc("/io.png", func(w http.ResponseWriter, r *http.Request) {
		plot(inputs[5], "IO", w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, html)
	})

	log.Printf("Listening on %s", ":3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
