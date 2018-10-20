// renderer provides a simple 3d renderer for simple objects onto Cube:Bit.
package renderer

import (
	"image/color"
	"math"

	"github.com/9600org/cubebit"
)

// Object is a thing to be rendered.
type Object interface {
	// At returns the colour at the given point in space.
	At(x, y, z float64) color.RGBA
}

// Sphere represents a sphere to be rendered
type Sphere struct {
	// CentreX, CentreY, and CentreZ specify the centre of the sphere in space.
	// The visible space is in the range [0..1]
	CentreX, CentreY, CentreZ float64
	// Radius is the radius of the sphere.
	Radius float64
	// CentreColour is the colour of the sphere at its centre.
	CentreColour color.RGBA
	// EdgeColour is the colour of the sphere at its edge.
	EdgeColour color.RGBA
}

// Renderer is a *very* simple renderer for objects on the Cube:Bit volume.
type Renderer struct {
	c *cubebit.Cubebit

	objects []Object
}

// New creates a new Renderer.
func New(c *cubebit.Cubebit) *Renderer {
	return &Renderer{c: c}
}

// Add adds an object to be rendered.
func (r *Renderer) Add(o Object) {
	r.objects = append(r.objects, o)
}

// gamma is a gamma correction table.
// TODO(al): Move this down to the cubebit package.
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

// Render renders the objects onto the Cube:Bit LEDs.
func (r *Renderer) Render() {
	sx, sz, sy := r.c.Bounds()
	sxf := float64(sx - 1)
	syf := float64(sy - 1)
	szf := float64(sz - 1)
	for z := 0; z < sz; z++ {
		for y := 0; y < sy; y++ {
			for x := 0; x < sx; x++ {
				var lr, lg, lb, _ uint32
				n := uint32(0)
				for _, o := range r.objects {
					r, g, b, _ := o.At(float64(x)/sxf, float64(y)/syf, float64(z)/szf).RGBA()
					lr += r >> 8
					lg += g >> 8
					lb += b >> 8
					n++
				}
				r.c.Set(x, y, z, color.RGBA{gamma[uint8(lr/n)], gamma[uint8(lg/n)], gamma[uint8(lb/n)], 255})
			}
		}
	}
	r.c.Render()
}

// blend returns a colour between a and b, according to the ratio specified.
func blend(a, b color.Color, blend float64) color.RGBA {
	inv := float64(0) // float64(1)-blend
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	or, og, ob, _ := uint32(float64(ar)*blend+float64(br)*inv),
		uint32(float64(ag)*blend+float64(bg)*inv),
		uint32(float64(ab)*blend+float64(bb)*inv),
		uint32(float64(aa)*blend+float64(ba)*inv)
	return color.RGBA{uint8(or >> 8), uint8(og >> 8), uint8(ob >> 8), 255} //uint8(oa>>8)}
}

// At implements Object.At.
func (s *Sphere) At(x, y, z float64) color.RGBA {
	dx := (s.CentreX - x)
	dy := (s.CentreY - y)
	dz := (s.CentreZ - z)
	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if dist > s.Radius {
		return color.RGBA{0, 0, 0, 0}
	}
	return blend(s.CentreColour, s.EdgeColour, float64(1)-((dist*dist)/s.Radius))
}
