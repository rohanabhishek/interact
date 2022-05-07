// This file consists of the details of WebServer
package web

import (
	"fmt"
	room "interact/server/room"
	rest "interact/server/services/rest"
	"net/http"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type WebServer struct {
	addr          *string
	serverMux     *mux.Router
	roomInstances map[string]*room.RoomInstance
}

func NewWebServer(addr string) *WebServer {
	webServer := new(WebServer)
	webServer.addr = &addr
	webServer.serverMux = mux.NewRouter()
	webServer.Handlers()
	webServer.roomInstances = make(map[string]*room.RoomInstance)
	return webServer
}

func (server *WebServer) Handlers() {
	server.serverMux.HandleFunc("/createEvent", func(w http.ResponseWriter, r *http.Request) {
		// sample way to send the socketInstance, roomInstance
		roomId, roomInstance, err := server.NewRoomInstance()
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		response := rest.CreateInstanceResponse{
			RoomId: roomId,
			Error:  errMsg,
		}

		//add to room instnace id mapping
		server.roomInstances[roomId] = roomInstance

		rest.CreateInstanceHandler(w, r, roomInstance, response)
	}).Methods("POST")

	server.serverMux.HandleFunc("/{roomId}/joinEvent", func(w http.ResponseWriter, r *http.Request) {

		room, ok := server.getRoomInstance(w, r)

		if !ok {
			return
		}

		rest.JoinEventHandler(w, r, room)

	}).Methods("POST")

	// TODO: Add HTTP Method, schemes
	server.serverMux.HandleFunc("/{roomId}/sendResponse/{clientId}", func(w http.ResponseWriter, r *http.Request) {
		room, ok := server.getRoomInstance(w, r)

		if !ok {
			return
		}

		//TODO: Validate client id
		clientId := mux.Vars(r)["clientId"]

		rest.ClientsResponseHandler(w, r, room)

		//Register client to LiveResultsSocket
		success := room.LiveResultsHandler.RegisterClient(clientId)

		if !success {
			glog.Error("Client found but not registered in Results handler how??")
		}
	}).Methods("POST")

	server.serverMux.HandleFunc("/{roomId}/addLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		room, ok := server.getRoomInstance(w, r)

		if !ok {
			return
		}

		rest.AddLiveQuestionHandler(w, r, room)
	}).Methods("POST")

	server.serverMux.HandleFunc("/{roomId}/fetchCurrentState", func(w http.ResponseWriter, r *http.Request) {
		room, ok := server.getRoomInstance(w, r)

		if !ok {
			return
		}
		rest.FetchCurrentStateHandler(w, r, room)
	}).Methods("GET")

	server.serverMux.HandleFunc("/{roomId}/fetchLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		room, ok := server.getRoomInstance(w, r)

		if !ok {
			return
		}
		rest.FetchLiveQuestionHandler(w, r, room)
	}).Methods("GET")

	server.serverMux.HandleFunc("/{roomId}/endEvent", func(w http.ResponseWriter, r *http.Request) {
		room, ok := server.getRoomInstance(w, r)

		if !ok {
			return
		}

		roomId := room.GetRoomId()

		rest.EndEventHandler(w, r, room)

		delete(server.roomInstances, roomId)

	}).Methods("POST")

	server.serverMux.HandleFunc("/{roomId}/nextLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		room, ok := server.getRoomInstance(w, r)

		if !ok {
			return
		}
		rest.MoveToNextQuestionHandler(w, r, room)
	}).Methods("POST")

	server.serverMux.HandleFunc("/{roomId}/liveResults/{clientId}", func(w http.ResponseWriter, r *http.Request) {
		room, ok := server.getRoomInstance(w, r)

		if !ok {
			return
		}

		//TODO: Add validation of clients

		clientId := mux.Vars(r)["clientId"]

		client, err := room.LiveResultsHandler.ServeWebsocket(w, r, clientId)

		if err != nil {
			//TODO: send error response
		}

		/**
			If there is previous socket and somehow we lost connection, we need to close goroutine
			of the previous one
		**/
		if prevClient, ok := room.LiveResultsHandler.ClientsMapping[clientId]; ok {
			//If previous socket is registered, we need to add this socket

			if room.LiveResultsHandler.IsClientRegistered(prevClient) {
				//register the new socket
				room.LiveResultsHandler.Register <- client

				//unregister the old socket
				room.LiveResultsHandler.Unregister <- prevClient
			}

			//close the go routine
			go func() {
				if prevClient != nil {
					prevClient.Close <- true
				}
				prevClient = nil
			}()
		}

		//Replace or Add new client
		room.LiveResultsHandler.ClientsMapping[clientId] = client
	})

	server.serverMux.HandleFunc("/{roomId}/liveQuestion/{clientId}", func(w http.ResponseWriter, r *http.Request) {
		room, ok := server.getRoomInstance(w, r)

		if !ok {
			return
		}

		//TODO: Add validation of clients

		clientId := mux.Vars(r)["clientId"]

		client, err := room.LiveQuestionHandler.ServeWebsocket(w, r, clientId)

		if err != nil {
			//TODO: send error response
		}

		/**
			If there is previous socket and somehow we lost connection, we need to close goroutine
			of the previous one
		**/
		if prevClient, ok := room.LiveQuestionHandler.ClientsMapping[clientId]; ok {
			//If previous socket is registered, we need to add this socket

			if room.LiveQuestionHandler.IsClientRegistered(prevClient) {
				//register the new socket
				room.LiveQuestionHandler.Register <- client

				//unregister the old socket
				room.LiveQuestionHandler.Unregister <- prevClient
			}

			//close the go routine
			go func() {
				if prevClient != nil {
					prevClient.Close <- true
				}
				prevClient = nil
			}()
		}
		// else {
		// 	//If client is not present which means he joins first time, so register him
		// 	go func() {
		// 		room.LiveQuestionHandler.Register <- client
		// 	}()
		// }

		//Replace or Add new client
		room.LiveQuestionHandler.ClientsMapping[clientId] = client
	})
}

func (server *WebServer) getRoomInstance(w http.ResponseWriter, r *http.Request) (room *room.RoomInstance, ok bool) {
	vars := mux.Vars(r)

	roomId := vars["roomId"]

	roomInstance, ok := server.roomInstances[roomId]

	if !ok {
		//Bad Request error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": Invalid room id %q}`, roomId)
		return nil, false
	}

	return roomInstance, true
}

// TODO: add shutdown of server
func (server *WebServer) Run() {
	glog.Info("Server listening on ", *server.addr)

	// Here we can use ListenAndServeTLS also
	glog.Fatal(http.ListenAndServe(*server.addr, server.serverMux))
}

func (server *WebServer) NewRoomInstance() (string, *room.RoomInstance, error) {
	roomId := uuid.NewString()

	glog.Info("New room instance is created", roomId)

	return roomId, room.NewRoomInstance(roomId), nil
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
