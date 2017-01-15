package main

import (
	"github.com/ziutek/de"
	"github.com/ziutek/kview"
	"github.com/ziutek/matrix"
	"golang.org/x/net/websocket"
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
	listenOn = "127.0.0.1:8080"
)

func html(w http.ResponseWriter, r *http.Request) {
	page.Exec(w, Ctx{ListenOn: listenOn})
}

type cost struct{}

func (_ cost) Cost(m matrix.Dense) float64 {
	return Cost(m.Elems())
}

func newcost() de.Cost {
	return cost{}
}

func data(w *websocket.Conn) {
	defer w.Close()

	min := matrix.AsDense(1, 4, Min())
	max := matrix.AsDense(1, 4, Max())
	m := de.NewMinimizer(newcost, 20, min, max)
	points := make([][2]int, len(m.Pop))
	for {
		minCost, maxCost := m.NextGen()
		// Transform agents in 4D space to points on 2D image
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
		// Check the end condition
		sum := math.Abs(minCost) + math.Abs(maxCost)
		diff := math.Abs(maxCost - minCost)
		if diff/(sum+2*math.SmallestNonzeroFloat64) < 1e-4 {
			return
		}
		// Slow down calculations for better presentation
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	page = kview.New("page.kt")
	http.HandleFunc("/", html)
	http.Handle("/img", Img)
	http.Handle("/data", websocket.Handler(data))
	log.Printf("Web server (at %s) is ready.\n", listenOn)
	http.ListenAndServe(listenOn, nil)
}
