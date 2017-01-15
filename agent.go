package de

import (
	"github.com/ziutek/matrix"
	"math/rand"
	"time"
)

type args struct {
	a, b, c matrix.Dense // three vectors for crossover
	f, cr   float64      // differential weight and crossover probability
}

type Agent struct {
	X, x matrix.Dense
	in   chan args    // to send three vectors for crossovera
	out  chan float64 // to obtain actual cost value for this agent
	rnd  *rand.Rand
	cost Cost
}

func newAgent(min, width matrix.Dense, cost Cost) *Agent {
	a := new(Agent)
	a.X = matrix.MakeDense(min.Size())
	a.x = a.X.AsRow() // crossover operates on vectorized matrix
	a.in = make(chan args, 1)
	a.out = make(chan float64, 1)
	a.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	a.cost = cost
	// Place this agent in random place on the initial area.
	for i := 0; i < a.x.NumRow(); i++ {
		for k := 0; k < a.x.NumCol(); k++ {
			a.x.Set(i, k, a.rnd.Float64())
		}
	}
	a.X.ArrMulBy(width)
	a.X.AddTo(min, 1)
	// Run crossover loop
	go a.crossoverLoop()
	return a
}

// Returns random cut points for crossover
func (a *Agent) cutpoints(cr float64) (start, stop int) {
	n := a.x.NumCol()
	l := 1
	for l < n && a.rnd.Float64() < cr {
		l++
	}
	// now l <= n
	start = a.rnd.Intn(n)
	stop = (start + l) % (n + 1)
	return
}

func perturb(u matrix.Dense, in args, start, stop int) {
	v := u.Cols(start, stop)
	v.Sub(in.b.Cols(start, stop), in.c.Cols(start, stop), in.f)
	v.AddTo(in.a.Cols(start, stop), 1)
}

// Crossover loop
func (a *Agent) crossoverLoop() {
	U := matrix.MakeDense(a.X.Size()) // place for mutated matrix
	u := U.AsRow()                    // we operate on vectorized matrices
	n := u.NumCol()
	costX := a.cost(a.X)

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
		if costU := a.cost(U); costU <= costX {
			a.X, U = U, a.X
			a.x, u = u, a.x
			costX = costU
		}
		// Return new cost
		a.out <- costX
	}
}
