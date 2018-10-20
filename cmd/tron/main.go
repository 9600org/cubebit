package main

import (
	"container/list"
	"flag"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/9600org/cubebit"
	"github.com/9600org/go-rpi-ws281x"
)

var (
	tailLen   = flag.Int("tail_len", 20, "length of tails")
	numCycles = flag.Int("cycles", 2, "Number of light cycles")
	frameDelay = flag.Duration("frame_delay", 50*time.Millisecond, "Inter-frame delay")
)

const (
	north = iota
	east  = iota
	up    = iota
	down  = iota
	west  = iota
	south = iota
)

type point struct {
	x, y, z int
}

func (p *point) add(dir int) point {
	r := *p
	switch dir {
	case north:
		r.y++
	case south:
		r.y--
	case east:
		r.x++
	case west:
		r.x--
	case up:
		r.z++
	case down:
		r.z--
	}
	return r
}

type cycle struct {
	dir        int
	head, last point
	tail       *list.List

	headColour color.RGBA
	tailColour color.RGBA
}

func (c *cycle) move(v func(point) bool) {
	valids := map[int]bool{north: true, south: true, east: true, west: true, up: true, down: true}
	// remove reverse as an option:
	delete(valids, 6-c.dir)
	if len(valids) == 0 {
		// TODO: boom
		return
	}

	c.tail.PushFront(c.head)
	if c.tail.Len() > *tailLen {
		c.tail.Remove(c.tail.Back())
	}

	switch r := rand.Float32(); {
	case r < 0.6:
		newHead := c.head.add(c.dir)
		if v(newHead) {
			c.head = newHead
			break
		}
		fallthrough
	default:

		// pick a new direction at random
		for newDir, _ := range valids {
			newHead := c.head.add(newDir)
			if v(newHead) {
				c.dir = newDir
				c.head = newHead
				break
			}
		}
	}
}

func (c *cycle) render(cube *cubebit.Cubebit) {

	tc := c.tailColour
	for ne := c.tail.Front(); ne != nil; ne = ne.Next() {
		neck := ne.Value.(point)
		cube.Set(neck.x, neck.y, neck.z, tc)
		d := uint8(20)
		if tc.R > d {
			tc.R -= d
		}
		if tc.G > d {
			tc.G -= d
		}
		if tc.B > d {
			tc.B -= d
		}
	}

	cube.Set(c.head.x, c.head.y, c.head.z, c.headColour)

	bum := c.tail.Back().Value.(point)
	cube.Set(bum.x, bum.y, bum.z, color.RGBA{0, 0, 0, 255})

}

func valid(c *cubebit.Cubebit) func(p point) bool {
	bx, by, bz := c.Bounds()
	return func(p point) bool {
		inside := p.x >= 0 && p.x < bx &&
			p.y >= 0 && p.y < by &&
			p.z >= 0 && p.z < bz
		if !inside {
			return false
		}
		col := c.At(p.x, p.y, p.z)
		return col.R == 0 && col.G == 0 && col.B == 0
	}
}

func main() {
	flag.Parse()
	config := ws281x.DefaultConfig
	config.Brightness = 255
	cube, done, err := cubebit.New(&config, 5, 5, 5)
	if err != nil {
		log.Fatalf("Failed to create new cubebit: %v", err)
	}
	defer done()

	cycles := []*cycle{}
	for i := 1; i <= *numCycles; i++ {
		r, g, b := uint8(127), uint8(127), uint8(127)
		var r2, g2, b2 uint8
		if i&0x01 != 0 {
			r = 255
			r2 = 255
		}
		if i&0x02 != 0 {
			g = 255
			g2 = 255
		}
		if i&0x04 != 0 {
			b = 255
			b2 = 255
		}
		cycles = append(cycles, &cycle{
			headColour: color.RGBA{r, g, b, 255},
			tailColour: color.RGBA{r2, g2, b2, 255},
			head:       point{rand.Intn(5), rand.Intn(5), rand.Intn(5)},
			tail:       list.New(),
		})
	}

	for {
		for _, c := range cycles {
			c.move(valid(cube))
			c.render(cube)
		}
		cube.Render()
		time.Sleep(*frameDelay)
	}
}
