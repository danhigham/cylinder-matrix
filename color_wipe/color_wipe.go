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

package color_wipe

import (
	_ "image/png"
	"time"
)

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

type ColorWipe struct {
	ws wsEngine
}

func (cw *ColorWipe) Setup(ws wsEngine) error {
	cw.ws = ws
	return cw.ws.Init()
}

func (cw *ColorWipe) Display(color uint32) error {
	// bounds := image.Rectangle{
	// 	Min: image.Point{X: 0, Y: 0},
	// 	Max: image.Point{X: width, Y: height},
	// }
	// for y := 0; y < height; y++ {
	// 	for x := 0; x < width; x++ {
	// 		cw.ws.Leds(0)[utils.CoordinatesToIndex(bounds, x, y, height, false)] = color
	// 		cw.ws.Render()
	// 		fmt.Printf("X:%d,Y:%d,I:%d\n", x, y, utils.CoordinatesToIndex(bounds, x, y, height, false))
	// 		time.Sleep(1000 * time.Millisecond)
	// 	}
	// }

	for i := 0; i < 100; i++ {
		cw.ws.Leds(0)[i] = color
		cw.ws.Render()
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
