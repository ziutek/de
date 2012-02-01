package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
)

const (
	scale = 50 // pixels/unit
	minX = -5
	maxX = 5
	minY = -4
	maxY = 4
)

var (
	img  []byte
)

// This type describes the optimization problem for both: the optimizer and
// the vizualization framework.
type OP struct{}

// Function to minimize
func (o OP) Cost(x, y float64) float64 {
	x += -1.5
	y += -2.5
	r := math.Sqrt(x*x + y*y)
	return -math.Cos(math.Pi*r) / (1 + r)
}

// Image generator
func (o OP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(img)
}

func init() {
	// Generate image that shows the cost function in the area
	b := image.Rect(
		0, 0,
		int((maxX-minX)*scale), int((maxY-minY)*scale),
	)
	m := image.NewNRGBA(b)
	c := color.NRGBA{A: 255}
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			f := ctx.OP.Cost(
				minX+float64(x)/scale,
				minY+float64(y)/scale,
			)
			if f < 0 {
				c.G, c.B = 0, byte(-255*f)
			} else {
				c.G, c.B = byte(255*f), 0
			}
			m.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, m)
	img = buf.Bytes()
}
