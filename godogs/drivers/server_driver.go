package drivers

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var server *http.Server
var upgrader = websocket.Upgrader{}
var confirmError error
var respond string
var receiveMessages []Message

// Message struct
type Message struct {
	value string
}

// StartServer a http websocket server
func StartServer(requestChannel chan string, responseChannel chan string) {
	router := http.NewServeMux()
	router.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		echo(w, r, requestChannel, responseChannel)
	})

	server = &http.Server{
		Addr:    addr,
		Handler: wrapHandler(router),
	}
	go server.ListenAndServe()
}

func echo(w http.ResponseWriter, r *http.Request, requestChannel chan string, responseChannel chan string) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info("upgrade:", err)
		return
	}
	defer c.Close()
	_, m, _ := c.ReadMessage()
	receiveMessages = append(receiveMessages, Message{value: string(m)})
	respond = "avadakedavra"
	if string(m) == "alohomora" {
		respond = "imperio"
	}
	confirmError = c.WriteMessage(websocket.TextMessage, []byte(respond))
	responseChannel <- respond

	newMessage := <-requestChannel
	c.WriteMessage(websocket.TextMessage, []byte(newMessage))
	_, res, _ := c.ReadMessage()
	responseChannel <- string(res)
}

func wrapHandler(f http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f.ServeHTTP(w, r)
	}
}

// CloseServer function
func CloseServer() {
	server.Close()
}

// ExistsPattern function
func ExistsPattern(pattern string) bool {
	for _, b := range receiveMessages {
		if b.value == pattern {
			return true
		}
	}
	return false
}

// NotExistsPattern function
func NotExistsPattern(pattern string) bool {
	return !ExistsPattern(pattern)
}

// HasError function
func HasError() bool {
	return confirmError != nil
}
