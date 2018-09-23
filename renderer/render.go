package renderer

import (
	"image/color"
	"math"
	// "fmt"

	"github.com/9600org/cubebit"
)

type Object interface {
	At(x, y, z float64) color.RGBA
}

type Sphere struct {
	CentreX, CentreY, CentreZ, Radius      float64
	CentreColor, EdgeColour color.RGBA
}

type Renderer struct {
	c *cubebit.Cubebit

	objects []Object
}

func New(c *cubebit.Cubebit) *Renderer {
	return &Renderer{c: c}
}

func (r *Renderer) Add(o Object) {
				r.objects = append(r.objects, o)
}

func (r *Renderer) Render() {
	sx, sz, sy := r.c.Bounds()
	sxf := float64(sx-1)
	syf := float64(sy-1)
	szf := float64(sz-1)
	for z := 0; z < sz; z++ {
		for y := 0; y < sy; y++ {
			for x := 0; x < sx; x++ {
				var lr, lg, lb, _ uint32
				n := uint32(0)
				for _, o := range r.objects {
					r, g, b, _ := o.At(float64(x)/sxf, float64(y)/syf, float64(z)/szf).RGBA()
					lr += r>>8
					lg += g>>8
					lb += b>>8
					n++
				}
				/*
				if lr <= 10 {lr = 10}
				if lg <= 10 {lg = 10}
				if lb <= 10 {lb = 10}
				*/
				r.c.Set(x, y, z, color.RGBA{uint8(lr/n), uint8(lg/n), uint8(lb/n), 255})
			}
		}
	}
	r.c.Render()
}

func blend(a, b color.Color, blend float64) color.RGBA {
				inv := float64(0) // float64(1)-blend
				ar, ag, ab, aa := a.RGBA()
				br, bg, bb, ba := b.RGBA()
				or, og, ob, _ := uint32(float64(ar)*blend+float64(br)*inv),
												  uint32(float64(ag)*blend+float64(bg)*inv),
													uint32(float64(ab)*blend+float64(bb)*inv),
													uint32(float64(aa)*blend+float64(ba)*inv)
				return color.RGBA{uint8(or>>8), uint8(og>>8), uint8(ob>>8), 255}//uint8(oa>>8)}
}

func (s *Sphere) At(x, y, z float64) color.RGBA {
				dx := (s.CentreX-x)
				dy := (s.CentreY-y)
				dz := (s.CentreZ-z)
				dist := math.Sqrt(dx*dx+dy*dy+dz*dz)
				if dist > s.Radius {
								return color.RGBA{0, 0, 0, 0}
				}
				return blend(s.CentreColor, s.EdgeColour, float64(1)-((dist*dist)/s.Radius))
}
