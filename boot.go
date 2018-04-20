package main

import (
	"flag"
	"log"

	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/app"
)

var addr = flag.String("addr", "127.0.0.1:7700", "imperium server address")

type websocketDialerShim struct {
	*websocket.Dialer
}

func (s websocketDialerShim) Dial(urlStr string) (app.WsConn, error) {
	conn, _, err := s.Dialer.Dial(urlStr, nil)
	return conn, err
}

func main() {
	flag.Parse()
	err := app.Start(*addr, new(websocketDialerShim))

	if err != nil {
		log.Println(err.Error())
	}
}
