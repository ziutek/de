// Differential evolution optimization library
//
// See http://en.wikipedia.org/wiki/Differential_evolution 
package de

import (
	"github.com/ziutek/matrix"
)

type abc struct {
	a, b, c *matrix.Dense
}

type agent struct {
	x   *matrix.Dense
	in  chan abc     // to send three vectors for crossover
	out chan float64 // to obtain actual cost value for this agent
	p   *Params
}

func newAgent(min, width *matrix.Dense, p *Params) *agent {
	a := new(agent)
	a.x = matrix.DenseZero(min.Size())
	a.in = make(chan abc, 1)
	a.out = make(chan float64, 1)
	a.p = p
	// Place this agent in random place on the initial area
	a.x.Rand(0, 1)
	a.x.MulBy(width)
	a.x.AddTo(min, 1)
	return a
}

// Crossover loop
func (a *agent) crLoop() {
	p := a.p
	u := matrix.DenseZero(a.x.Size())
	for in := range a.in {
		// crx contains three vectors from three random agents
		u.Sub(in.b, in.c, p.f)
	}
}

type Params struct {
	CR   float64                     // crossover probability (default 0.9)
	cost func(*matrix.Dense) float64 // function to minimize
	f    float64                     // differential weight
}

// Minimizes cost function by evoluting the population of multiple agents
type Minimizer struct {
	Params Params
	p      []*agent // population of agents
}

// Creates new minimizer
// cost - function to minimize
// n - number of entities in population
// min, max - area for initial population
func New(cost func(*matrix.Dense) float64, n int, min, max *matrix.Dense) *Minimizer {
	m := new(Minimizer)
	m.p = make([]*agent, n)
	m.Params.CR = 0.9
	m.Params.cost = cost
	// Initialization of population
	max.AddTo(min, -1) // calculate width of initial area
	for i := range m.p {
		m.p[i] = newAgent(min, max, &m.Params)
	}
	return m
}
