package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/app"
)

var addr = flag.String("addr", "127.0.0.1:7700", "imperium server address")

type wsShim struct {
	*websocket.Dialer
}

func (s wsShim) Dial(urlStr string, requestHeader http.Header) (app.WsConn, *http.Response, error) {
	return s.Dialer.Dial(urlStr, requestHeader)
}

func main() {
	flag.Parse()
	err := app.Start(*addr, wsShim{})

	if err != nil {
		log.Println(err.Error())
	}
}
