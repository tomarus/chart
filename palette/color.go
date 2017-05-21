package palette

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"regexp"
)

var rx3 = regexp.MustCompile(`(?i)^#?([0-9a-f]{3})$`)
var rx6 = regexp.MustCompile(`(?i)^#?([0-9a-f]{6})$`)

// hex2col converts a hex color with format [0-9a-f]{6} to a color.Color
func hex2col(in string) (color.Color, error) {
	dst := make([]byte, 3)
	_, err := hex.Decode(dst, []byte(in))
	return color.RGBA{dst[0], dst[1], dst[2], 255}, err
}

// hex2hsl converts a hex color with format [0-9a-f]{6} to a HSL color.
func hex2hsl(in string) (HSL, error) {
	dst := make([]byte, 3)
	_, err := hex.Decode(dst, []byte(in))
	return HSL{float64(360 / 255 * dst[0]), EightTo1(dst[1]), EightTo1(dst[2]), 1.}, err
}

// col2hex converts a color.RGBA to a hex string.
func col2hex(in *color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", in.R, in.G, in.B)
}

// ParseColor parses a color with format #abc or #aabbcc
func ParseColor(in string) (color.Color, error) {
	if rx3.MatchString(in) {
		s := rx3.FindStringSubmatch(in)
		ib := fmt.Sprintf("%c%c%c%c%c%c", s[1][0], s[1][0], s[1][1], s[1][1], s[1][2], s[1][2]) // meh
		return hex2col(ib)
	} else if rx6.MatchString(in) {
		s := rx6.FindStringSubmatch(in)
		return hex2col(s[1])
	} else {
		return nil, fmt.Errorf("can't parse color %s", in)
	}
}
