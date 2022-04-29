package socket

// This file consists of using websocket package

import (
	data "interact/server/config/data"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

//TODO: Remove this, temporary websocket handler
func WebSocketHandler(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) *websocket.Conn {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during connection upgradation:", err)
		return nil
	}
	defer ws.Close()

	reader(ws)
	return ws
}

func reader(conn *websocket.Conn) {
	for {

		str := "{\"question\":\"WhoistheCaptainofIndianCricketTeam\",\"results\":[{\"option\":\"kohli\",\"percentage\":20},{\"option\":\"Rohit\",\"percentage\":50},{\"option\":\"Pant\",\"percentage\":30}]}"
		if err := conn.WriteMessage(1, []byte(str)); err != nil {
			log.Println(err)
			return
		}

	}
}

func ServeWebsocket(ch *ClientHandler, w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client connected")

	client := &Client{ClientHandler: ch, conn: conn, send: make(chan []byte, 256)}
	client.ClientHandler.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeLiveResults()
}
