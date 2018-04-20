package receiver

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/connection"
)

type JobProcessor interface {
	Execute() (bool, error)
}

// Start function
func Start(conn connection.WsConn, jobProcessor JobProcessor) {
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
