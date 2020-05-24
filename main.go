package main

import (
	"flag"

	"github.com/danhigham/cylinder-matrix/marquee"
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
	rabbitHost := flag.String("rabbit-host", "127.0.0.1", "The RabbitMQ host to connect to for messages")
	rabbitPort := flag.Int("rabbit-port", 5672, "The RabbitMQ port to connect to for messages")
	rabbitQueue := flag.String("rabbit-queue", "marquee.messages", "The RabbitMQ queue to subscribe to for messages")
	rabbitUser := flag.String("rabbit-username", "guest", "The RabbitMQ user")
	rabbitPassword := flag.String("rabbit-password", "guest", "The RabbitMQ password")

	flag.Parse()

	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCounts

	dev, err := ws2811.MakeWS2811(&opt)
	check(err)

	marquee := &marquee.Marquee{}
	// colorWipe := &color_wipe.ColorWipe{}

	check(marquee.Setup(dev, *rabbitHost, *rabbitQueue, *rabbitUser, *rabbitPassword, *rabbitPort))
	marquee.WaitForMessages()
	// check(colorWipe.Setup(dev))
	defer dev.Fini()

	//	marquee.Display(args[0], 60, uint32(0xffffff))
	// colorWipe.Display(uint32(0xff0000))
	// colorWipe.Display(uint32(0x00ff00))
	// colorWipe.Display(uint32(0x0000ff))
	// colorWipe.Display(uint32(0x000000))

}
