package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/ziutek/de"
	"github.com/ziutek/kview"
	"github.com/ziutek/matrix"
	"log"
	"math"
	"net"
	"net/http"
	"syscall"
	"time"
)

type Ctx struct {
	ListenOn string
	OP       OP
}

var (
	page kview.View
	ctx  = Ctx{ListenOn: "127.0.0.1:8080"}
)

func html(w http.ResponseWriter, r *http.Request) {
	page.Exec(w, ctx)
}

func cost(m *matrix.Dense) float64 {
	return ctx.OP.Cost(m.Elems())
}

func data(w *websocket.Conn) {
	defer w.Close()

	min := matrix.NewDense(1, 2, 2, ctx.OP.Min()...)
	max := matrix.NewDense(1, 2, 2, ctx.OP.Max()...)
	p := matrix.DenseZero(1, 2)
	scale := ctx.OP.Scale()
	m := de.NewMinimizer(cost, 10, min, max)
	points := make([]struct{ X, Y int }, len(m.Pop))

	for {
		minCost, maxCost := m.NextGen()
		sum := math.Abs(minCost) + math.Abs(maxCost)
		diff := math.Abs(maxCost - minCost)
		if diff/(sum+2*math.SmallestNonzeroFloat64) < 1e-3 {
			return
		}

		for i, a := range m.Pop {
			// Calculate the coordinates of a point on the image
			p.Sub(a.X, min, scale)
			points[i].X = int(p.Get(0, 0))
			points[i].Y = int(p.Get(0, 1))
		}
		// Send points to the browser side application		
		if err := websocket.JSON.Send(w, points); err != nil {
			e, ok := err.(*net.OpError)
			if !ok || e.Err != syscall.EPIPE && e.Err != syscall.ECONNRESET {
				log.Print(err)
			}
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	page = kview.New("page.kt")
	http.HandleFunc("/", html)
	http.Handle("/img", ctx.OP)
	http.Handle("/data", websocket.Handler(data))
	http.ListenAndServe(ctx.ListenOn, nil)
}
