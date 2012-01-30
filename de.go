// Differential evolution optimization library
//
// See http://en.wikipedia.org/wiki/Differential_evolution 
package de

import (
	"github.com/ziutek/blas"
	"math"
)

type MinMax struct {
	Min, Max float64
}

type abc struct {
	a, b, c []float64
}

type entity struct {
	x, u []float64
	in   chan abc
	out  chan float64
}

type Minimizer struct {
	p  []entity // population
	f  float64  // mutation factor
	cr float64  // crossover probability
	n  int      // dimensionality of the problem

	cost func([]float64) float64 // function to minimize
}

func NewMinimizer(mm []MinMax, cr float64) *Minimizer {

}
