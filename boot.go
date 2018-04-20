package main

import (
	"flag"
	"log"

	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/connection"
	"github.com/segurosfalabella/imperium-worker/executer"
	"github.com/segurosfalabella/imperium-worker/receiver"
)

var addr = flag.String("addr", "127.0.0.1:7700", "imperium server address")

type websocketDialerShim struct {
	*websocket.Dialer
}

func (s websocketDialerShim) Dial(urlStr string) (connection.WsConn, error) {
	conn, _, err := s.Dialer.Dial(urlStr, nil)
	return conn, err
}

func main() {
	flag.Parse()
	conn, err := connection.Create(*addr, new(websocketDialerShim))

	if err != nil {
		log.Println(err.Error())
	}

	job := new(executer.Job)
	receiver.Start(conn, job)
}
