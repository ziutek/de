package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
)

const (
	width   = 60
	height  = 50
	minX    = -6
	maxX    = 6
	minY    = -5
	maxY    = 5

	minU    = -2.5
	maxU    = 2.5
	minV    = -1.5
	maxV    = 1.5


	mapW = int((maxX - minX) * ScaleXY)
	mapH = int((maxY - minY) * ScaleXY)
	cols = int((maxU - minU) * ScaleUV)
	rows = int((maxV - minV) * ScaleUV)
)

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

// Converts 4D point to 2D image point
func D4toD2(p []float64) (int, int) {
	x, y, u, v := p[0], p[1], p[2], p[3]
	k := int((u - minU) * ScaleUV) // column
	i := int((v - minV) * ScaleUV) // row
	return k*mapW + int(x*ScaleXY), i*mapH + int(y*ScaleXY)
}

// Converts 2D image point to 4D point
func D2toD4(i, k int) []float64 {
	x := minX + float64(i%mapW)/ScaleXY
	y := minY + float64(k%mapH)/ScaleXY
	u := minU + float64(i/mapW)/ScaleUV
	v := minV + float64(k/mapH)/ScaleUV
	return []float64{x, y, u, v}
}

// Generates image for specified u and v
func img(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	// Generate image that shows the cost function in the area
	b := image.Rect(
		0, 0,
		mapW*cols, mapH*rows,
	)
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
	png.Encode(w, m)
}
