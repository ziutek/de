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
	width  = 120
	height = 100
	minX   = -6
	maxX   = 6
	minY   = -5
	maxY   = 5

	cols = 9
	rows = 7
	minU = -3
	maxU = 3
	minV = -2.5
	maxV = 2.5

	dX     = maxX - minX
	dY     = maxY - minY
	dU     = maxU - minU
	dV     = maxV - minV
	scaleX = float64(dX) / float64(width)
	scaleY = float64(dY) / float64(height)
	scaleU = float64(dU) / float64(cols)
	scaleV = float64(dV) / float64(rows)
)

// Converts 4D point to 2D image point
func D4toD2(p []float64) (int, int) {
	x, y, u, v := p[0], p[1], p[2], p[3]
	c := int(cols * (u - minU) / dU)
	r := int(rows * (v - minV) / dV)
	// Check bounds
	if c < 0 || c >= cols || r < 0 || r >= rows {
		return -1, -1
	}
	i := int(width * (x - minX) / dX)
	k := int(height * (y - minY) / dY)
	// Check bounds
	if i < 0 || i >= width || k < 0 || k >= height {
		return -1, -1
	}
	return c*width + i, r*height + k
}

// Converts 2D image point to 4D point (notice that D2toD4(D4toD2(p)) != p)
func D2toD4(i, k int) []float64 {
	if i < 0 || k < 0 {
		nan := math.NaN()
		return []float64{nan, nan, nan, nan}
	}
	x := minX + (float64(i%width)+0.5)*scaleX
	y := minY + (float64(k%height)+0.5)*scaleY
	u := minU + (float64(i/width)+0.5)*scaleU
	v := minV + (float64(k/height)+0.5)*scaleV
	return []float64{x, y, u, v}
}

// Function to minimize
func Cost(p []float64) float64 {
	x := p[0]
	y := p[1]
	u := p[2]
	v := p[3]
	r := math.Sqrt(x*x + y*y + u*u + v*v)
	return -math.Cos(math.Pi*r) / (1 + r)
}

func Min() []float64 {
	return []float64{minU, minV, minX, minY}
}

func Max() []float64 {
	return []float64{maxU, maxV, maxX, maxY}
}

type pngImage []byte

// Sends image via HTTP
func (i pngImage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(i)
}

var Img pngImage

func init() {
	// Generate image that shows the cost function in 4d area
	b := image.Rect(0, 0, width*cols, height*rows)
	m := image.NewNRGBA(b)
	c := color.NRGBA{A: 255}
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			f := Cost(D2toD4(x, y))
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
	Img = buf.Bytes()
}
