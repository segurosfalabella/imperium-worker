package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/connection"
)

// WebsocketDialerShim type
type WebsocketDialerShim struct {
	*websocket.Dialer
}

// Dial method
func (ws *WebsocketDialerShim) Dial(urlStr string) (connection.WsConn, error) {
	conn, _, err := ws.Dialer.Dial(urlStr, nil)
	return conn, err
}

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

// Server var
var Server *http.Server

func main() {
}
