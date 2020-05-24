// Copyright 2018 Jacques Supcik / HEIA-FR
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package marquee

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/danhigham/cylinder-matrix/utils"
	"github.com/disintegration/imaging"
	"github.com/streadway/amqp"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

var charDict = map[byte][]int{
	65:  []int{0, 5},   //A
	66:  []int{6, 5},   //B
	67:  []int{12, 5},  //C
	68:  []int{18, 5},  //D
	69:  []int{24, 5},  //E
	70:  []int{30, 5},  //F
	71:  []int{36, 5},  //G
	72:  []int{42, 5},  //H
	73:  []int{48, 1},  //I
	74:  []int{50, 5},  //J
	75:  []int{56, 4},  //K
	76:  []int{61, 5},  //L
	77:  []int{67, 5},  //M
	78:  []int{73, 5},  //N
	79:  []int{79, 5},  //O
	80:  []int{85, 5},  //P
	81:  []int{91, 5},  //Q
	82:  []int{97, 5},  //R
	83:  []int{103, 5}, //S
	84:  []int{109, 5}, //T
	85:  []int{115, 5}, //U
	86:  []int{121, 5}, //V
	87:  []int{127, 5}, //W
	88:  []int{133, 5}, //X
	89:  []int{139, 5}, //Y
	90:  []int{145, 5}, //Z
	97:  []int{0, 5},   //a
	98:  []int{6, 5},   //b
	99:  []int{12, 5},  //c
	100: []int{18, 5},  //d
	101: []int{24, 5},  //e
	102: []int{30, 5},  //f
	103: []int{36, 5},  //g
	104: []int{42, 5},  //h
	105: []int{48, 1},  //i
	106: []int{50, 5},  //j
	107: []int{56, 4},  //k
	108: []int{61, 5},  //l
	109: []int{67, 5},  //m
	110: []int{73, 5},  //n
	111: []int{79, 5},  //o
	112: []int{85, 5},  //p
	113: []int{91, 5},  //q
	114: []int{97, 5},  //r
	115: []int{103, 5}, //s
	116: []int{109, 5}, //t
	117: []int{115, 5}, //u
	118: []int{121, 5}, //v
	119: []int{127, 5}, //w
	120: []int{133, 5}, //x
	121: []int{139, 5}, //y
	122: []int{145, 5}, //z
	49:  []int{151, 2}, //1
	50:  []int{154, 4}, //2
	51:  []int{159, 4}, //3
	52:  []int{164, 4}, //4
	53:  []int{169, 4}, //5
	54:  []int{174, 4}, //6
	55:  []int{179, 3}, //7
	56:  []int{183, 5}, //8
	57:  []int{189, 5}, //9
	48:  []int{195, 5}, //0
	33:  []int{201, 1}, //!
	35:  []int{203, 5}, //#
	36:  []int{209, 5}, //$
	37:  []int{215, 5}, //%
	94:  []int{221, 3}, //^
	38:  []int{225, 4}, //&
	42:  []int{230, 3}, //*
	40:  []int{234, 2}, //(
	41:  []int{237, 2}, //)
	95:  []int{240, 5}, //_
	43:  []int{246, 3}, //+
	61:  []int{250, 3}, //=
	45:  []int{254, 3}, //-
	47:  []int{258, 5}, ///
	92:  []int{264, 5}, //\
	126: []int{270, 5}, //~
	60:  []int{276, 2}, //<
	62:  []int{279, 2}, //>
	44:  []int{282, 2}, //,
	46:  []int{285, 1}, //.
}

const (
	brightness = 255
	width      = 20
	height     = 5
	ledCounts  = width * height
	maxCount   = 50
	sleepTime  = 200
)

type wsEngine interface {
	Init() error
	Render() error
	Wait() error
	Fini()
	Leds(channel int) []uint32
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type MarqueeMessage struct {
	Message string `json:"message"`
	Color   string `json:"color"`
}

type Marquee struct {
	charmap    map[byte]image.Image
	ws         wsEngine
	rabbitHost string
	rabbitPort int
	queue      string
	username   string
	password   string
}

func (m *Marquee) Setup(ws wsEngine, rabbitHost, queue, username, password string, rabbitPort int) error {

	m.ws = ws
	m.rabbitHost = rabbitHost
	m.rabbitPort = rabbitPort
	m.queue = queue
	m.username = username
	m.password = password
	m.charmap = make(map[byte]image.Image)

	charMapFile, err := os.Open("charmap.png")
	if err != nil {
		log.Fatal(err)
	}
	charmapPNG, err := png.Decode(charMapFile)
	if err != nil {
		log.Fatal(err)
	}

	for k, c := range charDict {
		cut := image.Rectangle{
			Min: image.Point{X: c[0], Y: 0},
			Max: image.Point{X: c[0] + c[1], Y: 5},
		}

		m.charmap[k] = imaging.Clone(charmapPNG.(SubImager).SubImage(cut))

	}

	return m.ws.Init()
}

func (m *Marquee) WaitForMessages() error {

	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", m.username, m.password, m.rabbitHost, m.rabbitPort)

	conn, err := amqp.Dial(uri)
	check(err)
	defer conn.Close()

	ch, err := conn.Channel()
	check(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		m.queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	check(err)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	check(err)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg MarqueeMessage
			err := json.Unmarshal(d.Body, &msg)
			check(err)
			rgb, err := utils.ParseHexColor(msg.Color)
			c := utils.RGBToColor(uint32(rgb.R)*255, uint32(rgb.G)*255, uint32(rgb.B)*255)
			m.Display([]byte(msg.Message), 60, c)
		}
	}()

	<-forever
	return nil
}

func (m *Marquee) Display(message []byte, delay int, color uint32) error {

	compositeWidth := 40

	for _, c := range message {
		if c == 32 {
			compositeWidth += 3
		}
		if char, ok := m.charmap[c]; ok {
			compositeWidth += char.Bounds().Max.X + 1
		}
	}

	r := image.Rectangle{image.Point{0, 0}, image.Point{compositeWidth, 5}}
	rgba := image.NewRGBA(r)

	currentPos := 20

	for _, c := range []byte(message) {
		if c == 32 {
			currentPos += 3
		}
		if char, ok := m.charmap[c]; ok {
			bounds := char.Bounds()
			r2 := image.Rectangle{image.Point{currentPos, 0}, image.Point{currentPos + bounds.Max.X, 5}}

			draw.Draw(rgba, r2, char, image.Point{0, 0}, draw.Src)
			currentPos += bounds.Max.X + 1
		}
	}

	// out, err := os.Create("./output.png")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer out.Close()

	// png.Encode(out, rgba)

	//do something here

	bounds := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height},
	}

	for offset := 0; offset < (rgba.Bounds().Max.X - width); offset++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				var drawColor uint32
				r, _, _, _ := rgba.At(x+offset, y).RGBA()
				if r == 65535 {
					drawColor = color
				} else {
					drawColor = uint32(0x000000)
				}
				m.ws.Leds(0)[utils.CoordinatesToIndex(bounds, x, y, height, true)] = drawColor
			}
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
		m.ws.Render()
	}

	return nil
}
