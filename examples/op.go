package main

import (
	"bytes"
	"github.com/ziutek/de"
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
)

const (
	scale = 50 // pixels/unit
)

var (
	area = []de.Range{{-4, 4}, {-4, 4}}
	img  []byte
)

// This type describes the optimization problem for both: the optimizer and
// the vizualization framework.
type OP struct{}

// Returns an initial area on which we search minimum of cost function
func (o OP) Area() []de.Range {
	return area
}

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
	x0 := area[0].Min
	y0 := area[1].Min
	b := image.Rect(
		0, 0,
		int((area[0].Max-x0)*scale), int((area[1].Max-y0)*scale),
	)
	m := image.NewNRGBA(b)
	c := color.NRGBA{A: 255}
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			f := ctx.OP.Cost(
				x0+float64(x)/scale,
				y0+float64(y)/scale,
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
