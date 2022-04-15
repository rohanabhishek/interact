// This file consists of the details of WebServer
package web

import (
	"errors"
	"net/http"

	data "interact/server/config/data"
	rest "interact/server/services/rest"

	"github.com/golang/glog"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
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
		roomId, err := server.NewRoomInstance()
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		response := rest.CreateInstanceResponse{
			RoomId: roomId,
			Error:  errMsg,
		}
		rest.CreateInstanceHandler(w, r, server.roomInstance, response)
	})
	// TODO: Add HTTP Method, schemes
	server.serverMux.HandleFunc("/{roomId}/sendResponse", func(w http.ResponseWriter, r *http.Request) {
		rest.ClientsResponseHandler(w, r, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/addLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		rest.AddLiveQuestionHandler(w, r, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/fetchLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		rest.FetchLiveQuestionHandler(w, r, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/endEvent", func(w http.ResponseWriter, r *http.Request) {
		rest.EndEventHandler(w, r, server.roomInstance)
	})

	server.serverMux.HandleFunc("/{roomId}/nextLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		rest.MoveToNextQuestionHandler(w, r, server.roomInstance)
	})
}

func (server *WebServer) Run() {
	glog.Info("Server listening on", *server.addr)
	// Here we can use ListenAndServeTLS also
	glog.Fatal(http.ListenAndServe(*server.addr, server.serverMux))
}

func (server *WebServer) NewRoomInstance() (string, error) {
	if server.roomInstance != nil {
		glog.Error("server roomInstance is not nil, Attempt to overwrite it")
		return "", errors.New("server roomInstance is not nil, Attempt to overwrite it")
	}
	server.roomInstance = data.NewRoomInstance()
	return server.roomInstance.GetRoomId(), nil
}
