// Differential evolution optimization library
// See http://en.wikipedia.org/wiki/Differential_evolution 
package de

import (
	"github.com/ziutek/matrix"
	"math"
	"math/rand"
	"time"
)

// Type of function to minimize
type Cost func(*matrix.Dense) float64

// Minimizes cost function by evoluting the population of multiple agents
type Minimizer struct {
	Pop []*Agent // population of agents
	CR  float64  // crossover probability (default 0.9)

	rnd *rand.Rand
}

// Creates new minimizer
// cost - function to minimize
// n - number of entities in population
// min, max - area for initial population
func New(cost Cost, n int, min, max *matrix.Dense) *Minimizer {
	if n < 4 {
		panic("population too small")
	}
	m := new(Minimizer)
	m.CR = 0.9
	m.Pop = make([]*Agent, n)
	// Initialization of population
	max.AddTo(min, -1) // calculate width of initial area
	for i := range m.Pop {
		m.Pop[i] = newAgent(min, max, cost)
	}
	m.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	return m
}

// Calculate next generation. Returns id and cost of the best agent. 
func (m *Minimizer) NextGen() (int, float64) {
	bestId := 0
	bestCost := math.MaxFloat64
	// Perform crossover
	for i, x := range m.Pop {
		f := 0.5 + m.rnd.Float64()*0.5
		a, b, c := m.abc(i)
		x.in <- args{a.x, b.x, c.x, f, m.CR}
	}
	// Get results
	for i, x := range m.Pop {
		c := <-x.out
		if c < bestCost {
			bestCost = c
			bestId = i
		}
	}
	return bestId, bestCost
}

// Stops all gorutines and invalidates m
func (m *Minimizer) Delete() {
	for _, x := range m.Pop {
		x.in <- args{}
	}
	m.Pop = nil
}

// Returns three random agents different from Pop[i]
func (m *Minimizer) abc(i int) (a, b, c *Agent) {
	n := len(m.Pop)
	j := m.rnd.Intn(n - 1)
	if j >= i {
		j++
	}
	k := m.rnd.Intn(n - 2)
	if k >= i {
		k++
	}
	if k >= j {
		k++
	}
	l := m.rnd.Intn(n - 3)
	if l >= i {
		l++
	}
	if l >= j {
		l++
	}
	if l >= k {
		l++
	}
	return m.Pop[j], m.Pop[k], m.Pop[l]
}
