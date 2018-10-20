package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/9600org/cubebit"
	"github.com/9600org/go-rpi-ws281x"
)

func main() {
	config := ws281x.DefaultConfig
	config.Brightness = 255
	c, done, err := cubebit.New(&config, 5, 5, 5)
	if err != nil {
		log.Fatalf("Failed to create new cubebit: %v", err)
	}
	defer done()

	col := color.RGBA{255, 0, 0, 255}
	for i := 0; ; i++ {
		for z := 0; z < 5; z++ {
			for y := 0; y < 5; y++ {
				for x := 0; x < 5; x++ {
					col.R = uint8(127 + 127*math.Sin(float64(i*2+x+y*5+z*25)/float64(5)))
					col.G = uint8(127 + 127*math.Sin(float64(i+x+y*5+z*25)/float64(10)))
					col.B = uint8(127 + 127*math.Sin(float64(i*3+x+y*5+z*25)/float64(15)))
					c.Set(x, y, z, col)
				}
			}
		}
		c.Render()
		time.Sleep(20 * time.Millisecond)
	}
}
