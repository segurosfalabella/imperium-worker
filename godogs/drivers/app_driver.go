package drivers

import (
	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/connection"
	"github.com/segurosfalabella/imperium-worker/executer"
	"github.com/segurosfalabella/imperium-worker/receiver"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type websocketDialerShim struct {
	*websocket.Dialer
}

func (s websocketDialerShim) Dial(urlStr string) (connection.WsConn, error) {
	conn, _, err := s.Dialer.Dial(urlStr, nil)
	return conn, err
}

const addr = "127.0.0.1:7700"

// RunApp function
func RunApp() {
	conn, err := connection.Create(addr, new(websocketDialerShim))
	if err != nil {
		log.Error(err.Error())
	}

	job := new(executer.Job)
	go receiver.Start(conn, job)
}
