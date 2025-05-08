package server

import (
	"fmt"
	"math"
)

func hslToRgb(h, s, l float64) (uint8, uint8, uint8) {
	var r, g, b float64

	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	switch {
	case h >= 0 && h < 60:
		r, g, b = c, x, 0
	case h >= 60 && h < 120:
		r, g, b = x, c, 0
	case h >= 120 && h < 180:
		r, g, b = 0, c, x
	case h >= 180 && h < 240:
		r, g, b = 0, x, c
	case h >= 240 && h < 300:
		r, g, b = x, 0, c
	default: // h >= 300 && h < 360
		r, g, b = c, 0, x
	}

	return uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255)
}

// GeneratePlayerColorStrings generates a slice of n distinct hexadecimal color strings.
func GeneratePlayerColorStrings(n int) []string {
	if n <= 0 {
		return []string{}
	}

	colors := make([]string, n)
	saturation := 0.7 // Adjust for desired color intensity
	lightness := 0.5  // Adjust for desired brightness

	for i := range n {
		hue := float64(i) * (360.0 / float64(n))
		r, g, b := hslToRgb(hue, saturation, lightness)
		colors[i] = fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}

	return colors
}
