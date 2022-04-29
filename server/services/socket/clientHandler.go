package socket

type ClientHandler struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	//incoming messages from clients
	Broadcast chan []byte
}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (ch *ClientHandler) Run() {
	for {
		select {

		//register the client
		case client := <-ch.register:
			ch.clients[client] = true

		//unregsiter the client
		case client := <-ch.unregister:
			if _, ok := ch.clients[client]; ok {
				delete(ch.clients, client)
				close(client.send)
			}

		//broadcast message to all registered clients
		case message := <-ch.Broadcast:
			for client := range ch.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(ch.clients, client)
				}
			}
		}
	}
}
