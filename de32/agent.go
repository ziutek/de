package de32

import (
	"github.com/ziutek/matrix/matrix32"
	"math/rand"
	"time"
)

type args struct {
	a, b, c matrix32.Dense // three vectors for crossover
	f, cr   float32      // differential weight and crossover probability
}

type Agent struct {
	x   matrix32.Dense
	in  chan args    // to send three vectors for crossovera
	out chan float32 // to obtain actual cost value for this agent
	rnd *rand.Rand
}

func newAgent(min, max []float32, newcost func() Cost) *Agent {
	a := new(Agent)
	a.x = matrix32.MakeDense(1, len(min))
	a.in = make(chan args, 1)
	a.out = make(chan float32, 1)
	a.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	// Place this agent in random place on the initial area.
	x := a.x.Elems()
	for i, m := range min {
		x[i] = m + (max[i]-m)*a.rnd.Float32()
	}
	// Run crossover loop
	go a.crossoverLoop(newcost)
	return a
}

func (a *Agent) X() []float32 {
	return a.x.Elems()
}

// Returns random cut points for crossover
func (a *Agent) cutpoints(cr float32) (start, stop int) {
	n := a.x.NumCol()
	l := 1
	for l < n && a.rnd.Float32() < cr {
		l++
	}
	// now l <= n
	start = a.rnd.Intn(n)
	stop = (start + l) % (n + 1)
	return
}

func perturb(u matrix32.Dense, in args, start, stop int) {
	v := u.Cols(start, stop)
	v.Sub(in.b.Cols(start, stop), in.c.Cols(start, stop), in.f)
	v.AddTo(in.a.Cols(start, stop), 1)
}

// Crossover loop
func (a *Agent) crossoverLoop(newcost func() Cost) {
	u := matrix32.MakeDense(a.x.Size())
	n := u.NumCol()
	cost := newcost()
	costX := cost.Cost(a.x.Elems())

	for in := range a.in {
		if !in.a.IsValid() {
			return
		}
		start, stop := a.cutpoints(in.cr)
		// Crossover:
		//  u = a + f * (b + c) for elements between start and stop
		//  u = x               for remaining elements
		if start < stop {
			u.Cols(0, start).Copy(a.x.Cols(0, start))
			perturb(u, in, start, stop)
			u.Cols(stop, n).Copy(a.x.Cols(stop, n))
		} else {
			stop++ // because it is modulo (n+1)
			perturb(u, in, 0, stop)
			u.Cols(stop, start).Copy(a.x.Cols(stop, start))
			perturb(u, in, start, n)
		}
		// Selection
		if costU := cost.Cost(u.Elems()); costU <= costX {
			a.x, u = u, a.x
			costX = costU
		}
		// Return new cost
		a.out <- costX
	}
}
