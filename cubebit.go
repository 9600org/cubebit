package cubebit

import (
	"image/color"

	"github.com/mcuadros/go-rpi-ws281x"
)

type Cubebit struct {
	canvas *ws281x.Canvas
	sizeX int
	sizeY int
	sizeZ int
}


func New(config *ws281x.HardwareConfig, sx, sy, sz int) (*Cubebit, func() error, error) {
	width := sx
	height := sy*sz
	c, err := ws281x.NewCanvas(width, height, config)
	if err != nil {
		return nil, nil, err
	}
	if err := c.Initialize(); err != nil {
		return nil, c.Close, err
	}
	return &Cubebit{canvas: c, sizeX: sx, sizeY: sy, sizeZ: sx}, c.Close, nil
}

func (c *Cubebit) Set(x, y, z int, col color.RGBA) {
	if z%2==1 {
		y, x = c.sizeY-x-1, c.sizeX-y-1
	}
	if y%2==1 {
		x = c.sizeX-x-1
	}
	page := c.sizeY*z
	c.canvas.Set(x, page+y, col)
}

func (c *Cubebit) Render() {
	c.canvas.Render()
}

func (c *Cubebit) Bounds() (int, int, int) {
	return c.sizeX, c.sizeY, c.sizeZ
}

