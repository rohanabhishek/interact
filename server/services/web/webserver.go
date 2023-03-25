// This file consists of the details of WebServer
package web

import (
	"encoding/json"
	room "interact/server/room"
	rest "interact/server/services/rest"
	"io"
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
		glog.V(2).Info("/createEvent", r)
		if r.Method == "POST" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				glog.Error("IO Request Body read failed", err)
				errorResponse := rest.ErrorResponse{Error: err.Error()}
				rest.SetResponseMetadata(w)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errorResponse)
				return
			}

			err = room.ValidateRoomUnMarshal(bodyBytes)
			if err != nil {
				errorResponse := rest.ErrorResponse{Error: err.Error()}
				rest.SetResponseMetadata(w)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errorResponse)
				return
			}

			// sample way to send the socketInstance, roomInstance
			roomId, roomInstance, err := server.NewRoomInstance()
			if err != nil {
				errorResponse := rest.ErrorResponse{Error: err.Error()}
				rest.SetResponseMetadata(w)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errorResponse)
				return
			}

			//add to room instance id mapping
			server.roomInstances[roomId] = roomInstance

			err = roomInstance.UnMarshal(bodyBytes)
			if err != nil {
				errorResponse := rest.ErrorResponse{Error: err.Error()}
				rest.SetResponseMetadata(w)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errorResponse)
				return
			}

			createEventResponse := rest.CreateInstanceResponse{
				RoomId: roomInstance.GetRoomId(),
				HostId: roomInstance.GetHostId(),
			}
			rest.SetResponseMetadata(w)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(createEventResponse)
		} else {
			rest.SetResponseMetadata(w)
			w.WriteHeader(http.StatusOK)
		}
	}).Methods("POST", "OPTIONS")

	server.serverMux.HandleFunc("/{roomId}/joinEvent", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {

			room, ok := server.getRoomInstance(w, r)

			if !ok {
				glog.Errorf("Room not found %s", mux.Vars(r)["roomId"])
				return
			}

			rest.JoinEventHandler(w, r, room)
		} else {
			rest.SetResponseMetadata(w)
			w.WriteHeader(http.StatusOK)
		}
	}).Methods("POST", "OPTIONS")

	// TODO: Add HTTP Method, schemes
	server.serverMux.HandleFunc("/{roomId}/sendResponse/{clientId}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {

			room, ok := server.getRoomInstance(w, r)

			if !ok {
				glog.Errorf("Room not found %s", mux.Vars(r)["roomId"])
				return
			}

			//TODO: Validate client id
			clientId := mux.Vars(r)["clientId"]

			ok = rest.ClientsResponseHandler(w, r, room)
			if !ok {
				glog.Error("ClientResponseHandler error")
				return
			}

			success := room.SocketHandler.RegisterClient(clientId)

			if !success {
				glog.Error("Client found but not registered in Results handler how??")
			}
		} else {
			rest.SetResponseMetadata(w)
			w.WriteHeader(http.StatusOK)
		}
	}).Methods("POST", "OPTIONS")

	server.serverMux.HandleFunc("/{roomId}/addLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {

			room, ok := server.getRoomInstance(w, r)

			if !ok {
				glog.Errorf("Room not found %s", mux.Vars(r)["roomId"])
				return
			}

			rest.AddLiveQuestionHandler(w, r, room)
		} else {
			rest.SetResponseMetadata(w)
			w.WriteHeader(http.StatusOK)
		}
	}).Methods("POST", "OPTIONS")

	server.serverMux.HandleFunc("/{roomId}/fetchCurrentState", func(w http.ResponseWriter, r *http.Request) {
		room, ok := server.getRoomInstance(w, r)

		if !ok {
			glog.Errorf("Room not found %s", mux.Vars(r)["roomId"])
			return
		}
		rest.FetchCurrentStateHandler(w, r, room)
	}).Methods("GET", "OPTIONS")

	server.serverMux.HandleFunc("/{roomId}/fetchLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		room, ok := server.getRoomInstance(w, r)

		if !ok {
			glog.Errorf("Room not found %s", mux.Vars(r)["roomId"])
			return
		}
		rest.FetchLiveQuestionHandler(w, r, room)
	}).Methods("GET", "OPTIONS")

	server.serverMux.HandleFunc("/{roomId}/endEvent", func(w http.ResponseWriter, r *http.Request) {
		if (r.Method) == "POST" {

			room, ok := server.getRoomInstance(w, r)

			if !ok {
				glog.Errorf("Room not found %s", mux.Vars(r)["roomId"])
				return
			}

			roomId := room.GetRoomId()

			rest.EndEventHandler(w, r, room)

			delete(server.roomInstances, roomId)

		} else {
			rest.SetResponseMetadata(w)
			w.WriteHeader(http.StatusOK)
		}
	}).Methods("POST", "OPTIONS")

	server.serverMux.HandleFunc("/{roomId}/nextLiveQuestion", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			room, ok := server.getRoomInstance(w, r)

			if !ok {
				glog.Errorf("Room not found %s", mux.Vars(r)["roomId"])
				return
			}
			rest.MoveToNextQuestionHandler(w, r, room)
		} else {
			rest.SetResponseMetadata(w)
			w.WriteHeader(http.StatusOK)
		}
	}).Methods("POST", "OPTIONS")

	server.serverMux.HandleFunc("/{roomId}/ws/{clientId}", func(w http.ResponseWriter, r *http.Request) {
		room, ok := server.getRoomInstance(w, r)

		if !ok {
			glog.Errorf("Room not found %s", mux.Vars(r)["roomId"])
			return
		}

		//TODO: Add validation of clients
		clientId := mux.Vars(r)["clientId"]

		client, err := room.SocketHandler.ServeWebsocket(w, r, clientId)

		if err != nil {
			//TODO: send error response
			glog.Error("error serving client responses id:", clientId)
		}

		/**
			If there is previous socket and somehow we lost connection, we need to close goroutine
			of the previous one
		**/
		if prevClient, ok := room.SocketHandler.ClientsMapping[clientId]; ok {
			//If previous socket is registered, we need to add this socket
			glog.Info("existing socket was found for the given client... replacing with the new one id:", clientId)

			if room.SocketHandler.IsClientRegistered(prevClient) {

				glog.Info("registering the new socket client id:", clientId)
				//register the new socket
				room.SocketHandler.Register <- client

				glog.Info("unregistering the old socket client id:", clientId)
				//unregister the old socket
				room.SocketHandler.Unregister <- prevClient
			}

			//close the go routine
			go func() {
				if prevClient != nil {
					glog.Info("closing the previous go routine")
					prevClient.Close <- true
				}
				prevClient = nil
			}()
		}

		//Replace or Add new client
		glog.Info("adding client to socket client handler mapping id:", clientId)
		room.SocketHandler.ClientsMapping[clientId] = client
	})
}

func (server *WebServer) getRoomInstance(w http.ResponseWriter, r *http.Request) (room *room.RoomInstance, ok bool) {
	vars := mux.Vars(r)

	roomId := vars["roomId"]

	roomInstance, ok := server.roomInstances[roomId]

	if !ok {
		//Bad Request error
		errorResponse := rest.ErrorResponse{Error: "Invalid room id: " + roomId}
		rest.SetResponseMetadata(w)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
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

	hostId := uuid.NewString()

	glog.Info("New room instance is created: ", roomId)

	glog.Info("Host Id: ", hostId)

	return roomId, room.NewRoomInstance(roomId, hostId), nil
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
