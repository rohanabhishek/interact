// This file consists of the details of WebServer
package web

import (
	"net/http"

	"github.com/golang/glog"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	data "interact/server/config/data"
	rest "interact/server/services/rest"
)

type WebServer struct {
	addr           *string
	serverMux      *mux.Router
	socketInstance *socketio.Server
	roomInstance   *data.RoomInstance
}

func NewWebServer(addr string) *WebServer {
	webServer := new(WebServer)
	webServer.addr = &addr
	webServer.socketInstance = NewWebSocket()
	//TODO: Add event handling of the socket instance using the EventHandling
	// func in websocket.go OR add them appropriately when needed
	webServer.serverMux = mux.NewRouter()
	webServer.Handlers()
	return webServer
}

func (server *WebServer) Handlers() {
	server.serverMux.HandleFunc("/createEvent", func(w http.ResponseWriter, r *http.Request) {
		// sample way to send the socketInstance, roomInstance
		rest.CreateInstanceHandler(w, r, server.socketInstance, server.roomInstance)
	})
	// TODO: Add HTTP Method, schemes
	server.serverMux.HandleFunc("/{roomId}/sendResponse/{clientId}", func(w http.ResponseWriter, r *http.Request) {
		rest.ClientsResponseHandler(w, r, server.socketInstance, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/addLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		rest.AddLiveQuestion(w, r, server.socketInstance, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/fetchLiveQuestion/{clientId}", func(w http.ResponseWriter, r *http.Request) {
		rest.FetchLiveQuestion(w, r, server.socketInstance, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/endEvent", func(w http.ResponseWriter, r *http.Request) {
		rest.EndEvent(w, r, server.socketInstance, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/nextLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		rest.MoveToNextQuestion(w, r, server.socketInstance, server.roomInstance)
	})
}

func (server *WebServer) Run() {
	glog.Info("Server listening on", *server.addr)
	// Here we can use ListenAndServeTLS also
	glog.Fatal(http.ListenAndServe(*server.addr, server.serverMux))
}
