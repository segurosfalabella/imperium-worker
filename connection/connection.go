package connection

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// WsConn interface
type WsConn interface {
	Close() error
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type dialer interface {
	Dial(urlStr string) (WsConn, error)
}

// Create method
func Create(address string, dialer dialer) (WsConn, error) {
	if address == "" {
		return nil, errors.New("missing server address")
	}

	if validateAddressFormat(address) {
		return nil, errors.New("server address invalid")
	}

	url := url.URL{Scheme: "ws", Host: address, Path: "/azkaban"}
	conn, error := dialer.Dial(url.String())
	if error != nil {
		return nil, error
	}

	log.Info("connected to:", address)
	return conn, nil
}

func validateAddressFormat(address string) bool {
	ipFormat := `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`
	hostnameFormat := `\w+((.\w+)+)?`
	expression := fmt.Sprintf(`(%s)|(%s):\d+`, ipFormat, hostnameFormat)
	match, err := regexp.MatchString(expression, address)
	return err != nil || !match
}
