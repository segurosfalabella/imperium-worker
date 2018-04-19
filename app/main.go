package app

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/gorilla/websocket"
)

// WsDial variable
var WsDial = websocket.DefaultDialer.Dial

// Start method
func Start(address string) error {
	if address == "" {
		return errors.New("missing server address")
	}

	if validateAddressFormat(address) {
		return errors.New("server address invalid")
	}

	url := url.URL{Scheme: "ws", Host: address, Path: "/echo"}
	_, _, error := WsDial(url.String(), nil)
	if error != nil {
		return error
	}

	// fmt.Println(connection.ReadMessage())

	return nil
}

func validateAddressFormat(address string) bool {
	ipFormat := `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`
	hostnameFormat := `\w+((.\w+)+)?`
	expression := fmt.Sprintf(`(%s)|(%s):\d+`, ipFormat, hostnameFormat)
	match, err := regexp.MatchString(expression, address)
	return err != nil || !match
}
