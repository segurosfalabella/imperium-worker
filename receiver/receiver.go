package receiver

import (
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/connection"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

const passwordForSend = "alohomora"
const passwordForValidate = "imperio"

// JobProcessor interface
type JobProcessor interface {
	Execute() error
}

// Start function
func Start(conn connection.WsConn, jobProcessor JobProcessor) {
	err := auth(conn)
	if err == nil {
		loop(conn, jobProcessor)
	}
}

func auth(conn connection.WsConn) error {
	conn.WriteMessage(websocket.TextMessage, []byte(passwordForSend))
	_, message, _ := conn.ReadMessage()
	if string(message) != passwordForValidate {
		return errors.New("server unknown")
	}
	return nil
}

func loop(conn connection.WsConn, jobProcessor JobProcessor) {
	for {
		messageType, message, _ := conn.ReadMessage()
		switch messageType {
		case websocket.TextMessage:
			json.Unmarshal(message, &jobProcessor)
		default:
		}

		jobProcessor.Execute()
		// TODO: Salir de una manera elegante.
		return
	}
}
