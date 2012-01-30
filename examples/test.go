package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/ziutek/kview"
	"log"
	"net"
	"net/http"
	"syscall"
	"time"
)

var (
	listenOn = "127.0.0.1:8080"
	page     kview.View
)

type Ctx struct {
	HostPort string
}

func show(w http.ResponseWriter, r *http.Request) {
	ctx := Ctx{
		HostPort: listenOn,
	}
	page.Exec(w, ctx)
}

func data(w *websocket.Conn) {
	for i := 0; ; i++ {
		if err := websocket.JSON.Send(w, i); err != nil {
			if oe, ok := err.(*net.OpError); !ok || oe.Err != syscall.EPIPE {
				log.Print(err)
			}
			return
		}
		time.Sleep(time.Second)
	}
}

func main() {
	page = kview.New("page.kt")
	http.HandleFunc("/", show)
	http.Handle("/data", websocket.Handler(data))
	http.ListenAndServe(listenOn, nil)
}
