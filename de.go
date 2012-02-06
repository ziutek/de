// Differential evolution optimization library
//
// See http://en.wikipedia.org/wiki/Differential_evolution 
package de

import (
	"github.com/ziutek/matrix"
	"math/rand"
	"time"
)

type abc struct {
	a, b, c *matrix.Dense
}

type agent struct {
	X      *matrix.Dense
	in     chan abc     // to send three vectors for crossover
	out    chan float64 // to obtain actual cost value for this agent
	rnd    *rand.Rand
	params *Params
}

func newAgent(min, width *matrix.Dense, p *Params) *agent {
	a := new(agent)
	a.X = matrix.DenseZero(min.Size()).Hvec() // we operate on vectors only
	a.in = make(chan abc, 1)
	a.out = make(chan float64, 1)
	a.params = p
	a.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	// Place this agent in random place on the initial area
	a.X.Rand(0, 1)
	a.X.MulBy(width)
	a.X.AddTo(min, 1)
	return a
}

// Returns random cut points for crossover
func (a *agent) cutpoints() (start, stop int) {
	n := len(a.X.Elems())
	l := 1
	for l < n && a.rnd.Float64() < a.params.CR {
		l++
	}
	start = a.rnd.Intn(n)
	stop = (start + l) % n
	return
}

func perturb(u *matrix.Dense, in abc, f float64, start, stop int) {
	v := u.Vslice(start, stop)
	v.Sub(in.b.Vslice(start, stop), in.c.Vslice(start, stop), f)
	v.AddTo(in.a.Vslice(start, stop), 1)
}

// Crossover loop
func (a *agent) crossoverLoop() {
	u := matrix.DenseZero(a.X.Size())
	n := u.Rows()
	p := a.params

	for in := range a.in {
		x := a.X
		start, stop := a.cutpoints()
		// Crossover:
		//  u = a + f * (b + c) for elements between start and stop
		//  u = x               for remaining elements
		if start < stop {
			u.Vslice(0, start).Copy(x.Vslice(0, start))
			perturb(u, in, p.f, start, stop)
			u.Vslice(stop, n).Copy(x.Vslice(stop, n))
		} else {
			perturb(u, in, p.f, 0, stop)
			u.Vslice(stop, start).Copy(x.Vslice(start, stop))
			perturb(u, in, p.f, start, n)
		}
		// x = u if cost(u) <= cost(x)
		/*if (... ) {
			a.X = u
			u = x
		}*/
	}
}

type Params struct {
	CR   float64                     // crossover probability (default 0.9)
	cost func(*matrix.Dense) float64 // function to minimize
	f    float64                     // differential weight
}

// Minimizes cost function by evoluting the population of multiple agents
type Minimizer struct {
	Params     Params
	cols, rows int
	p          []*agent // population of agents
	best       *matrix.Dense
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
	m.cols, m.rows = min.Size()
	// Initialization of population
	max.AddTo(min, -1) // calculate width of initial area
	for i := range m.p {
		m.p[i] = newAgent(min, max, &m.Params)
	}
	return m
}

// Returns best resul
func (m *Minimizer) Best() *matrix.Dense {
	return matrix.NewDense(m.rows, m.cols, m.cols, m.best.Elems())
}
