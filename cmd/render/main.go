package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/9600org/cubebit"
	"github.com/9600org/cubebit/renderer"
	"github.com/9600org/go-rpi-ws281x"
)

func main() {
	config := ws281x.DefaultConfig
	config.Pin = 18
	c, done, err := cubebit.New(&config, 5, 5, 5)
	if err != nil {
		log.Fatalf("Failed to create new cubebit: %v", err)
	}
	defer done()

	r := renderer.New(c)

	s := []*renderer.Sphere{
		&renderer.Sphere{0.5, 0.5, 0.5, 0.9, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 0, 255}},
		&renderer.Sphere{0.5, 0.5, 0.5, 0.9, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 0, 255}},
	}
	for _, o := range s {
		r.Add(o)
	}

	for t := float64(0.1); ; t += float64(0.1) {
		for i, o := range s {
			o.CentreX, o.CentreY, o.CentreZ = posAtT(t, i)
		}
		r.Render()
		time.Sleep(10 * time.Millisecond)
	}
}

func posAtT(ti float64, tweak int) (float64, float64, float64) {
	t := float64(tweak+1) / 5.0
	x := 0.5 + 0.5*math.Sin(ti*t)
	y := 0.5 + 0.5*math.Sin(ti*1.2*t)
	z := 0.5 + 0.5*math.Sin(ti*1.5*t)
	return x, y, z
}
