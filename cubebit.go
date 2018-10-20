// cubebit provides an API for doing stuff with the Cube:Bit from 4tronix.
package cubebit

import (
	"image/color"

	"github.com/9600org/go-rpi-ws281x"
)

// Cubebit represents an instance of the Cube:Bit hardware.
type Cubebit struct {
	canvas *ws281x.Canvas
	sizeX  int
	sizeY  int
	sizeZ  int
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
	height := sy * sz
	c, err := ws281x.NewCanvas(width, height, config)
	if err != nil {
		return nil, nil, err
	}
	if err := c.Initialize(); err != nil {
		return nil, c.Close, err
	}
	return &Cubebit{canvas: c, sizeX: sx, sizeY: sy, sizeZ: sx}, c.Close, nil
}

// gamma is a gamma correction table.
var gamma = []uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2,
	2, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 5, 5, 5,
	5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 9, 9, 9, 10,
	10, 10, 11, 11, 11, 12, 12, 13, 13, 13, 14, 14, 15, 15, 16, 16,
	17, 17, 18, 18, 19, 19, 20, 20, 21, 21, 22, 22, 23, 24, 24, 25,
	25, 26, 27, 27, 28, 29, 29, 30, 31, 32, 32, 33, 34, 35, 35, 36,
	37, 38, 39, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 50,
	51, 52, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 66, 67, 68,
	69, 70, 72, 73, 74, 75, 77, 78, 79, 81, 82, 83, 85, 86, 87, 89,
	90, 92, 93, 95, 96, 98, 99, 101, 102, 104, 105, 107, 109, 110, 112, 114,
	115, 117, 119, 120, 122, 124, 126, 127, 129, 131, 133, 135, 137, 138, 140, 142,
	144, 146, 148, 150, 152, 154, 156, 158, 160, 162, 164, 167, 169, 171, 173, 175,
	177, 180, 182, 184, 186, 189, 191, 193, 196, 198, 200, 203, 205, 208, 210, 213,
	215, 218, 220, 223, 225, 228, 231, 233, 236, 239, 241, 244, 247, 249, 252, 255}

func gammafy(col color.RGBA) color.RGBA {
	return color.RGBA{gamma[col.R], gamma[col.G], gamma[col.B], col.A}
}

// Set turns the pixel at the specified coordinates to the specified colour.
// Note that no LEDs will change colour until the application also calls the
// Render function, below.
func (c *Cubebit) Set(x, y, z int, col color.RGBA) {
	if z%2 == 1 {
		y, x = c.sizeY-x-1, c.sizeX-y-1
	}
	if y%2 == 1 {
		x = c.sizeX - x - 1
	}
	page := c.sizeY * z
	c.canvas.Set(x, page+y, gammafy(col))
}

func (c *Cubebit) At(x, y, z int) color.RGBA {
	if z%2 == 1 {
		y, x = c.sizeY-x-1, c.sizeX-y-1
	}
	if y%2 == 1 {
		x = c.sizeX - x - 1
	}
	page := c.sizeY * z
	return c.canvas.At(x, page+y).(color.RGBA)
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
