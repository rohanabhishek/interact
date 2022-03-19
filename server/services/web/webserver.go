package web

import (
	"net/http"

	"github.com/golang/glog"
	socketio "github.com/googollee/go-socket.io"
	restHandler "interact/server/services/rest"
)

type WebServer struct {
	addr            *string
	serverMux       *http.ServeMux
	socketInstances map[string]*socketio.Server
}

func NewWebServer(addr *string) *WebServer {
	webServer := new(WebServer)
	webServer.addr = addr
	webServer.serverMux = http.NewServeMux()
	return webServer
}

func (server *WebServer) Handlers() {
	server.serverMux.HandleFunc("/createEvent", func(w http.ResponseWriter, r *http.Request) {
		restHandler.CreateInstanceHandler(w, r, server.socketInstances)
	})
}

func (server *WebServer) Run() {
	glog.Info("Server listening on", *server.addr)
	// Here we can use ListenAndServeTLS also
	glog.Fatal(http.ListenAndServe(*server.addr, server.serverMux))
}
