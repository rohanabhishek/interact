package socket

// This file consists of using websocket package

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func (ch *ClientHandler) ServeWebsocket(w http.ResponseWriter, r *http.Request, clientId string) (c *Client, e error) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Client connected")

	client := &Client{clientId: clientId, ClientHandler: ch, conn: conn, send: make(chan []byte, 256), Close: make(chan bool)}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeLiveResults()

	return client, nil
}
