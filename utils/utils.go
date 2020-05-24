package utils

import (
	"fmt"
	"image"
	"image/color"
)

func CoordinatesToIndex(bounds image.Rectangle, x, y, height int, flipY bool) int {
	if flipY {
		y = (bounds.Max.Y - 1) - y
	}

	if x%2 == 0 {
		return (x-bounds.Min.X)*height + (y - bounds.Min.Y)
	}

	return (x-bounds.Min.X)*height + (height - 1) - (y - bounds.Min.Y)
}

func RGBToColor(r uint32, g uint32, b uint32) uint32 {
	return ((r>>8)&0xff)<<16 + ((g>>8)&0xff)<<8 + ((b >> 8) & 0xff)
}

func ParseHexColor(s string) (c color.RGBA64, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}
