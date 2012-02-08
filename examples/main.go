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
}

var (
	page     kview.View
	ListenOn = "127.0.0.1:8080"
)

func html(w http.ResponseWriter, r *http.Request) {
	page.Exec(w, Ctx{ListenOn: ListenOn})
}

func cost(m *matrix.Dense) float64 {
	return Cost(m.Elems())
}

func data(w *websocket.Conn) {
	defer w.Close()

	min := matrix.NewDense(1, 4, 4, Min()...)
	max := matrix.NewDense(1, 4, 4, Max()...)
	m := de.NewMinimizer(cost, 20, min, max)
	points := make([][2]int, len(m.Pop))
	for {
		minCost, maxCost := m.NextGen()
		sum := math.Abs(minCost) + math.Abs(maxCost)
		diff := math.Abs(maxCost - minCost)
		if diff/(sum+2*math.SmallestNonzeroFloat64) < 1e-3 {
			return
		}
		// Transform agents to points on image
		for i, a := range m.Pop {
			points[i][0], points[i][1] = D4toD2(a.X.Elems())
		}
		// Send points to the browser
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
	http.Handle("/img", Img)
	http.Handle("/data", websocket.Handler(data))
	http.ListenAndServe(ListenOn, nil)
}
