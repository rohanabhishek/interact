package socket

import (
	"github.com/gorilla/websocket"
)

type Client struct {

	//id of the client
	clientId string

	ClientHandler *ClientHandler

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	//close the client
	Close chan bool
}

func (c *Client) writeLiveResults() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//Since we make sure only one message is sent, else we need to handle multiple messages
			c.conn.WriteMessage(1, message)
		case <-c.Close:
			return
		}
	}
}
