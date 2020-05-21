package main

import (
	"github.com/danhigham/cylinder-matrix/color_wipe"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	brightness = 255
	width      = 20
	height     = 5
	ledCounts  = width * height
	maxCount   = 50
	sleepTime  = 200
)

func main() {
	// args := os.Args[1:]

	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCounts

	dev, err := ws2811.MakeWS2811(&opt)
	check(err)

	// marquee := &marquee.Marquee{}
	colorWipe := &color_wipe.ColorWipe{}

	// check(marquee.Setup(dev))
	check(colorWipe.Setup(dev))
	defer dev.Fini()

	// marquee.Display(args[0])
	colorWipe.Display(uint32(0xff0000))
	colorWipe.Display(uint32(0x000000))

	// for count := 0; count < maxCount; count++ {
	// 	inv.display()
	// 	inv.next()
	// 	time.Sleep(sleepTime * time.Millisecond)
	// }
}
