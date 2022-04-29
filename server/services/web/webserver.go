// This file consists of the details of WebServer
package web

import (
	data "interact/server/config/data"
	rest "interact/server/services/rest"
	socket "interact/server/services/socket"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

type WebServer struct {
	addr          *string
	serverMux     *mux.Router
	roomInstance  *data.RoomInstance
	ClientHandler *socket.ClientHandler
}

func NewWebServer(addr string) *WebServer {
	webServer := new(WebServer)
	webServer.addr = &addr
	webServer.ClientHandler = socket.NewClientHandler()
	go webServer.ClientHandler.Run()
	//TODO: Add event handling of the socket instance using the EventHandling
	// func in websocket.go OR add them appropriately when needed
	webServer.serverMux = mux.NewRouter()
	webServer.Handlers()
	return webServer
}

func (server *WebServer) Handlers() {
	server.serverMux.HandleFunc("/createEvent", func(w http.ResponseWriter, r *http.Request) {
		// sample way to send the socketInstance, roomInstance
		rest.CreateInstanceHandler(w, r, server.roomInstance)
	})
	// TODO: Add HTTP Method, schemes
	server.serverMux.HandleFunc("/{roomId}/sendResponse/{clientId}", func(w http.ResponseWriter, r *http.Request) {
		rest.ClientsResponseHandler(w, r, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/addLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		rest.AddLiveQuestion(w, r, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/fetchLiveQuestion/{clientId}", func(w http.ResponseWriter, r *http.Request) {
		rest.FetchLiveQuestion(w, r, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/endEvent", func(w http.ResponseWriter, r *http.Request) {
		rest.EndEvent(w, r, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/nextLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		rest.MoveToNextQuestion(w, r, server.roomInstance)
	})

	//TODO: Remove this
	server.serverMux.HandleFunc("/{roomId}/socket", func(w http.ResponseWriter, r *http.Request) {
		socket.WebSocketHandler(w, r, server.roomInstance)
	})

	//TODO: add roomId to url
	server.serverMux.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWebsocket(server.ClientHandler, w, r)
	})

}

func (server *WebServer) Run() {
	glog.Info("Server listening on", *server.addr)

	//TODO: Remove this and proper json handling
	str := "{\"question\":\"WhoistheCaptainofIndianCricketTeam\",\"results\":[{\"option\":\"kohli\",\"percentage\":20},{\"option\":\"Rohit\",\"percentage\":50},{\"option\":\"Pant\",\"percentage\":30}]}"

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				server.ClientHandler.Broadcast <- []byte(str)
			}
		}
	}()

	// Here we can use ListenAndServeTLS also
	glog.Fatal(http.ListenAndServe(*server.addr, server.serverMux))

}
