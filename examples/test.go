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

func data(w *websocket.Conn) {
	for i := 0; ; i++ {
		if err := websocket.JSON.Send(w, i); err != nil {
			if e, ok := err.(*net.OpError); !ok || e.Err != syscall.EPIPE {
				log.Print(err)
			}
			return
		}
		time.Sleep(time.Second)
	}
}

func main() {
	page = kview.New("page.kt")
	http.HandleFunc("/", html)
	http.Handle("/func.png", ctx.OP)
	http.Handle("/data", websocket.Handler(data))
	http.ListenAndServe(ctx.ListenOn, nil)
}
