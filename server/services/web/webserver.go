// This file consists of the details of WebServer
package web

import (
	"errors"
	room "interact/server/room"
	rest "interact/server/services/rest"
	socket "interact/server/services/socket"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

type WebServer struct {
	addr      *string
	serverMux *mux.Router
	//TODO: Map of room instances
	roomInstance *room.RoomInstance
}

func NewWebServer(addr string) *WebServer {
	webServer := new(WebServer)
	webServer.addr = &addr
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
	}).Methods("POST")

	server.serverMux.HandleFunc("/{roomId}/joinEvent", func(w http.ResponseWriter, r *http.Request) {
		rest.JoinEventHandler(w, r, server.roomInstance)
	}).Methods("POST")

	// TODO: Add HTTP Method, schemes
	server.serverMux.HandleFunc("/{roomId}/sendResponse", func(w http.ResponseWriter, r *http.Request) {
		rest.ClientsResponseHandler(w, r, server.roomInstance)
	}).Methods("POST")

	server.serverMux.HandleFunc("/{roomId}/addLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		rest.AddLiveQuestionHandler(w, r, server.roomInstance)
	}).Methods("POST")

	server.serverMux.HandleFunc("/{roomId}/fetchCurrentState", func(w http.ResponseWriter, r *http.Request) {
		rest.FetchCurrentStateHandler(w, r, server.roomInstance)
	}).Methods("GET")

	server.serverMux.HandleFunc("/{roomId}/fetchLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		rest.FetchLiveQuestionHandler(w, r, server.roomInstance)
	}).Methods("GET")

	server.serverMux.HandleFunc("/{roomId}/endEvent", func(w http.ResponseWriter, r *http.Request) {
		rest.EndEventHandler(w, r, server.roomInstance)
	}).Methods("POST")

	server.serverMux.HandleFunc("/{roomId}/nextLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		rest.MoveToNextQuestionHandler(w, r, server.roomInstance)
	}).Methods("POST")

	server.serverMux.HandleFunc("/{roomId}/liveResults", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWebsocket(server.roomInstance.LiveResultsHandler, w, r)
	}).Methods("POST")

	server.serverMux.HandleFunc("/{roomId}/liveQuestion", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWebsocket(server.roomInstance.LiveQuestionHandler, w, r)
	}).Methods("POST")

}

// TODO: add shutdown of server
func (server *WebServer) Run() {
	glog.Info("Server listening on ", *server.addr)

	//TODO: Remove this and proper json handling

	str := `{"question":"Who is the Captain of Indian Cricket Team","results":[{"option":"Kohli","percentage":20},{"option":"Rohit","percentage":50},{"option":"Pant","percentage":30}]}`
	strq := `{"question":"Who is the Captain of Indian Cricket Team","answers":["Kohli", "Rohit", "Pant"]}`

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if server.roomInstance != nil {
					server.roomInstance.LiveResultsHandler.Broadcast <- []byte(str)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				if server.roomInstance != nil {
					server.roomInstance.LiveQuestionHandler.Broadcast <- []byte(strq)
				}
			}
		}
	}()

	// Here we can use ListenAndServeTLS also
	glog.Fatal(http.ListenAndServe(*server.addr, server.serverMux))
}

func (server *WebServer) NewRoomInstance() (string, error) {
	if server.roomInstance != nil {
		glog.Error("server roomInstance is not nil, Attempt to overwrite it")
		return "", errors.New("server roomInstance is not nil, Attempt to overwrite it")
	}
	server.roomInstance = room.NewRoomInstance()
	return server.roomInstance.GetRoomId(), nil
}

/*
APIs to be invoked from client side
- After a client joins using the JoinEvent API, it then triggers
	FetchCurrentState API,
		- if its response is WAITING_ON_CLIENTS_FOR_RESPONSES
			then trigger the FetchCurrentQuestion API
		- else display on client side WAITING_ON_HOST_FOR_QUESTION
	- This approach could cause an issue, as the server could simultaneously invoke
		the other event after the FetchCurrentState.
	- So, will be better to use Socket to gather the state??
- At the start of the event, AddLiveQuestion is directly invoked
 After the event has started, to add a new question, MoveToNextQuestion,
 AddLiveQuestion are invoked by server
*/
