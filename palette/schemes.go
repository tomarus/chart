package palette

import "sync"

var defaultSchemes = map[string]map[string]string{
	"white": {
		"background": "#fff",
		"title":      "#000",
		"title2":     "#888",
		"grid":       "#0000007f",
		"grid2":      "#0000003f",
		"border":     "#666",
		"marker":     "#666",
		"select":     "#999",
		"area":       "#f2545b",
		"color1":     "#a93f55",
		"color2":     "#8c5e58",
		"color3":     "#19323c",
	},
	"black": {
		"background": "#000",
		"title":      "#fff",
		"title2":     "#aaa",
		"grid":       "#9999997f",
		"grid2":      "#5555557f",
		"border":     "#777",
		"marker":     "#fff",
		"select":     "#0f0",
		"area":       "#0e402d",
		"color1":     "#295135",
		"color2":     "#5a6650",
		"color3":     "#9fcc2e",
	},
	"pink": {
		"background": "#f2d1ba",
		"title":      "#5e2728",
		"title2":     "#7e4748",
		"grid":       "#5e2728",
		"grid2":      "#ae7778",
		"border":     "#5e2728",
		"marker":     "#7e4748",
		"select":     "#eee",
		"area":       "#f34093",
		"color1":     "#f78bd1",
		"color2":     "#d2082d",
		"color3":     "#df60a3",
	},
	"solarized": { // http://ethanschoonover.com/solarized
		"background": "#002b36",
		"title":      "#eee8d5",
		"title2":     "#eee8d5",
		"grid":       "#eee8d57f",
		"grid2":      "#8e88757f",
		"border":     "#657b83",
		"marker":     "#fdf6e3",
		"select":     "#fdf6e3",
		"area":       "#93a1a1",
		"color1":     "#657b83",
		"color2":     "#586e75",
		"color3":     "#eee8d5",
	},
}

var defaultLock sync.RWMutex

func getDefaultScheme(name string) (map[string]string, bool) {
	defaultLock.RLock()
	defer defaultLock.RUnlock()
	x, ok := defaultSchemes[name]
	return x, ok
}

// AddScheme adds a color scheme to the library so it can be used on a
// next call to NewPalette()
func AddScheme(name string, scheme map[string]string) {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	defaultSchemes[name] = scheme
}
