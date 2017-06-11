package palette

import "sync"

var defaultSchemes = map[string]map[string]string{
	"white": {
		"background": "#fff",
		"title":      "#660",
		"grid":       "#0000003f",
		"border":     "#666",
		"marker":     "#666",
		"select":     "#999",
		"area":       "#8abb31",
		"color1":     "#0000ff",
		"color2":     "#ff0000",
		"color3":     "#ff00ff",
	},
	"black": {
		"background": "#000",
		"title":      "#fff",
		"grid":       "#4040403f",
		"border":     "#ccc",
		"marker":     "#fff",
		"select":     "#0f0",
		"area":       "#8abb31",
		"color1":     "#4444ff",
		"color2":     "#ff4444",
		"color3":     "#ff44ff",
	},
	"pink": {
		"background": "#f2d1ba",
		"title":      "#5e2728",
		"grid":       "#5e2728",
		"border":     "#5e2728",
		"marker":     "#7e4748",
		"select":     "#eee",
		"area":       "#f34093",
		"color1":     "#f78bd1",
		"color2":     "#d2082d",
		"color3":     "#cf0063",
	},
	"solarized": { // http://ethanschoonover.com/solarized
		"background": "#002b36",
		"title":      "#eee8d5",
		"grid":       "#586e75",
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
