package socket

import "github.com/golang/glog"

type ClientHandler struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	//incoming messages from clients
	Broadcast chan []byte

	//close the handler
	Close chan bool

	//All clients mapping with thier clientId
	ClientsMapping map[string]*Client
}

func (ch *ClientHandler) IsClientRegistered(c *Client) bool {
	if val, err := ch.clients[c]; !err {
		return val
	}
	return false
}

func (ch *ClientHandler) RegisterClient(id string) bool {
	if client, ok := ch.ClientsMapping[id]; ok {

		glog.Info("Registering client id:", id)

		if client != nil {
			ch.clients[client] = true
			return true
		}
	}

	glog.Error("Unable to register client id:", id)
	return false
}

func (ch *ClientHandler) RegisterAllClients() {
	for _, client := range ch.ClientsMapping {
		if client != nil {
			ch.clients[client] = true
		}
	}
}

func (ch *ClientHandler) UnRegisterAllClients() {

	glog.Info("Unregistering all clients...")

	for k := range ch.clients {
		delete(ch.clients, k)
	}
}

func NewClientHandler() *ClientHandler {
	ch := &ClientHandler{
		Broadcast:      make(chan []byte, 256),
		Register:       make(chan *Client, 256),
		Unregister:     make(chan *Client, 256),
		clients:        make(map[*Client]bool),
		Close:          make(chan bool),
		ClientsMapping: make(map[string]*Client),
	}

	//run client handler
	go ch.run()

	return ch
}

func (ch *ClientHandler) run() {
	for {
		select {

		//register the client
		case client := <-ch.Register:
			ch.clients[client] = true

		//unregsiter the client
		case client := <-ch.Unregister:
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
		case <-ch.Close:
			for client := range ch.clients {
				select {
				case client.Close <- true:
				default:
					close(client.send)
					delete(ch.clients, client)
				}
			}
		}
	}
}
