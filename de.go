// Differential evolution optimization library
//
// See http://en.wikipedia.org/wiki/Differential_evolution 
package de

import (
)

type Range struct {
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

// Creates new minimizer
// ia - area for initial population
// cr - crossover probability
func NewMinimizer(ia []Range, cr float64) *Minimizer {
	return nil
}
