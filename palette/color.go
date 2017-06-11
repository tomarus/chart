package palette

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"regexp"
)

var rx3 = regexp.MustCompile(`(?i)^#?([0-9a-f]{3})$`)
var rx6 = regexp.MustCompile(`(?i)^#?([0-9a-f]{6})$`)
var rx8 = regexp.MustCompile(`(?i)^#?([0-9a-f]{8})$`)

// hex2col converts a hex color with format [0-9a-f]{6} to a color.Color
func hex2col(in string) (*color.RGBA, error) {
	dst := make([]byte, 4)
	_, err := hex.Decode(dst, []byte(in))
	a := dst[3]
	if a == 0 {
		a = 255
	}
	return &color.RGBA{dst[0], dst[1], dst[2], a}, err
}

// col2hex converts a color.RGBA to a hex string.
func col2hex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r/256, g/256, b/256)
}

// ParseColor parses a color with format #abc or #aabbcc
func ParseColor(in string) (*color.RGBA, error) {
	if rx8.MatchString(in) {
		s := rx8.FindStringSubmatch(in)
		return hex2col(s[1])
	} else if rx6.MatchString(in) {
		s := rx6.FindStringSubmatch(in)
		return hex2col(s[1])
	} else if rx3.MatchString(in) {
		s := rx3.FindStringSubmatch(in)
		ib := fmt.Sprintf("%c%c%c%c%c%c", s[1][0], s[1][0], s[1][1], s[1][1], s[1][2], s[1][2]) // meh
		return hex2col(ib)
	}
	return nil, fmt.Errorf("can't parse color %s", in)
}
