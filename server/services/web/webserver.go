// This file consists of the details of WebServer
package web

import (
	"net/http"

	"github.com/golang/glog"
	socketio "github.com/googollee/go-socket.io"
	data "interact/server/config/data"
	restHandler "interact/server/services/rest"
)

type WebServer struct {
	addr           *string
	serverMux      *http.ServeMux
	socketInstance *socketio.Server
	roomInstance   *data.RoomInstance
}

func NewWebServer(addr string) *WebServer {
	webServer := new(WebServer)
	webServer.addr = &addr
	webServer.socketInstance = NewWebSocket()
	//TODO: Add event handling of the socket instance using the EventHandling
	// func in websocket.go OR add them appropriately when needed
	webServer.serverMux = http.NewServeMux()
	webServer.Handlers()
	return webServer
}

func (server *WebServer) Handlers() {
	server.serverMux.HandleFunc("/createEvent", func(w http.ResponseWriter, r *http.Request) {
		// sample way to send the socketInstance, roomInstance
		restHandler.CreateInstanceHandler(w, r, server.socketInstance, server.roomInstance)
	})
}

func (server *WebServer) Run() {
	glog.Info("Server listening on", *server.addr)
	// Here we can use ListenAndServeTLS also
	glog.Fatal(http.ListenAndServe(*server.addr, server.serverMux))
}
