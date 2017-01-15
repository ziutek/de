// Package de implements Differential Evolution method of optimization.
// See http://en.wikipedia.org/wiki/Differential_evolution
package de

import (
	"github.com/ziutek/matrix"
	"log"
	"math"
	"math/rand"
	"time"
)

// Cost is an interface that an optimization problem should implement.
type Cost interface {
	Cost(matrix.Dense) float64
}

// Minimizes cost function by evoluting the population of multiple agents.
type Minimizer struct {
	Pop    []*Agent // population of agents
	CR     float64  // crossover probability (default 0.9)
	BestId int

	rnd *rand.Rand
}

// Creates new minimizer:
//	newcost - should return new value that satisfies Cost interface (it is called
//	concurently by population of agents during minmizer initialization),
//	n - number of entities in population,
//	min, max - area for initial population.
func NewMinimizer(newcost func() Cost, n int, min, max matrix.Dense) *Minimizer {
	if n < 4 {
		log.Panic("population too small: ", n)
	}
	m := new(Minimizer)
	m.CR = 0.9
	m.Pop = make([]*Agent, n)
	// Initialization of population
	max.AddTo(min, -1) // calculate width of initial area
	for i := range m.Pop {
		m.Pop[i] = newAgent(min, max, newcost)
	}
	m.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	return m
}

// Calculate next generation. Returns min and max cost in population.
func (m *Minimizer) NextGen() (minCost, maxCost float64) {
	minCost = math.MaxFloat64
	maxCost = -math.MaxFloat64

	// Perform crossover
	for i, x := range m.Pop {
		f := 0.5 + m.rnd.Float64()*0.5
		a, b, c := m.abc(i)
		x.in <- args{a.x, b.x, c.x, f, m.CR}
	}
	// Get results
	for i, x := range m.Pop {
		c := <-x.out
		if c < minCost {
			minCost = c
			m.BestId = i
		}
		if c > maxCost {
			maxCost = c
		}
	}
	return
}

// Stops all gorutines and invalidates m.
func (m *Minimizer) Delete() {
	for _, x := range m.Pop {
		x.in <- args{}
	}
	m.Pop = nil
}

// Returns three random agents different from Pop[i].
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
