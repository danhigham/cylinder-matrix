package utils

import "image"

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
