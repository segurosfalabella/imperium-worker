package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gorilla/websocket"
)

// WsConn interface
type WsConn interface {
	Close() error
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

// Dialer interface
type Dialer interface {
	Dial(urlStr string, requestHeader http.Header) (WsConn, *http.Response, error)
}

// Start method
func Start(address string, dialer Dialer) error {
	if address == "" {
		return errors.New("missing server address")
	}

	if validateAddressFormat(address) {
		return errors.New("server address invalid")
	}

	url := url.URL{Scheme: "ws", Host: address, Path: "/echo"}
	conn, _, error := dialer.Dial(url.String(), nil)
	if error != nil {
		return error
	}

	log.Println("connected to:", address)

	channel := make(chan string)
	go loop(conn, channel)

	conn.WriteMessage(websocket.TextMessage, []byte("e16b7b57-3eab-4866-805a-81ccc15a01ac"))

	message := <-channel
	log.Println("message receive:", message)
	return nil
}

func validateAddressFormat(address string) bool {
	ipFormat := `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`
	hostnameFormat := `\w+((.\w+)+)?`
	expression := fmt.Sprintf(`(%s)|(%s):\d+`, ipFormat, hostnameFormat)
	match, err := regexp.MatchString(expression, address)
	return err != nil || !match
}

func loop(conn WsConn, channel chan string) {
	if _, message, err := conn.ReadMessage(); err != nil {
		log.Println(err.Error())
	} else {
		channel <- string(message)
	}
}
