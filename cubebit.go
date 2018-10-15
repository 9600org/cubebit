// cubebit provides an API for doing stuff with the Cube:Bit from 4tronix.
package cubebit

import (
	"image/color"

	"github.com/9600org/go-rpi-ws281x"
)

// Cubebit represents an instance of the Cube:Bit hardware.
type Cubebit struct {
	canvas *ws281x.Canvas
	sizeX int
	sizeY int
	sizeZ int
}


// New creates a new Cubebit instance.
// config passes the details of the hardware to the underlying libws281x library.
// sx, sy, sz specify the dimentions of the LED space (e.g. 5, 5, 5 for the
// 5x5x5 model).
// It returns a Cubebit instance, and a function which should be called to
// release hardware // resources once the caller is finished with the Cubebit
// instance.
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

// Set turns the pixel at the specified coordinates to the specified colour.
// Note that no LEDs will change colour until the application also calls the
// Render function, below.
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

// Render sends the colour data previously associated with LEDs using the Set
// function to the LEDs themselves.
func (c *Cubebit) Render() {
	c.canvas.Render()
}

// Bounds returns the upper limits of the LED space.
func (c *Cubebit) Bounds() (int, int, int) {
	return c.sizeX, c.sizeY, c.sizeZ
}

