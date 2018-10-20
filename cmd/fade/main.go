package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/9600org/cubebit"
	"github.com/mcuadros/go-rpi-ws281x"
)

func main() {
	config := ws281x.DefaultConfig
	config.Pin = 18
	c, done, err := cubebit.New(&config, 5, 5, 5)
	if err != nil {
		log.Fatalf("Failed to create new cubebit: %v", err)
	}
	defer done()

	for i := 0; ; i++ {
		p := float64(i)
		col := color.RGBA{uint8(127 + 127*math.Sin(p/100.0)),
			0, //uint8(127+63*math.Sin(p/200.0)),
			0, //uint8(127+63*math.Sin(p/300.0)),
			255}
		for z := 0; z < 5; z++ {
			for y := 0; y < 5; y++ {
				for x := 0; x < 5; x++ {
					c.Set(x, y, z, col)
				}
			}
		}
		c.Render()
		time.Sleep(20 * time.Millisecond)
	}
}
