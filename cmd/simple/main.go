package main

import (
	"log"
	"image/color"
	"time"

	"github.com/mcuadros/go-rpi-ws281x"
	"github.com/9600org/cubebit"
)

func main() {
	config := ws281x.DefaultConfig
	config.Pin = 18
	c, done, err := cubebit.New(&config, 5, 5, 5)
	if err != nil {
		log.Fatalf("Failed to create new cubebit: %v", err)
	}
	defer done()

	col := color.RGBA{255, 0, 0, 255}
	for z := 0; z < 5; z++ {
		for y := 0; y < 5; y++ {
			for x:= 0; x < 5; x++ {
				c.Set(x, y, z, col)
				c.Render()
			time.Sleep(200*time.Millisecond)
			}
		}
	}
}
