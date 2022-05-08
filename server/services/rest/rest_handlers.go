// This file consists of the handlers used by the Server
package rest

import (
	// "fmt"

	room "interact/server/room"
	"net/http"

	"encoding/json"

	"github.com/golang/glog"
	"github.com/google/uuid"

	// "github.com/gorilla/mux"
	"io"
)

func CreateInstanceHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance,
	roomInstanceResponse CreateInstanceResponse) {

	/*
		Usage of r:
		r.Method  // request method
		r.URL     // request URL
		r.Header  // request headers
		r.Body    // request body
		https://pkg.go.dev/net/http#Request
	*/
	glog.V(2).Info("CreateInstanceHandler: ", r)
	glog.V(2).Info("CreateInstanceHandler Body: ", r.Body)

	// Sample usecase to display text on webpage
	// vars := mux.Vars(r)
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div><div>%v</div>", "Interact",
	// 	"Application", vars)
	if roomInstanceResponse.Error != "" {
		json.NewEncoder(w).Encode(roomInstanceResponse)
		// resp, _ := json.Marshal(roomInstanceResponse)
		// w.Write(resp)
		return
	}

	// Pre-processing of the request body
	// bodyBytes, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	glog.Error("IO Request Body read failed", err)
	// 	roomInstanceResponse.Error = "IO Request Body read failed" + err.Error()
	// 	json.NewEncoder(w).Encode(roomInstanceResponse)
	// 	return
	// }

	// err = room.SetRoomConfig(bodyBytes)
	// if err != nil {
	// 	roomInstanceResponse.Error = err.Error()
	// }

	// TODO: set the status of the APIs appropriately in case of errors
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roomInstanceResponse)
	/*
		Two ways to write the response
		1. The currently used method in server\services\rest\rest_handlers.go and decoding in server
			server\test\server_handlers_test.go
		2.
			a. Encoding
				responseBytes, err := json.Marshal(response)
				w.Write(responseBytes)
			b. Decoding
				defer response.Body.Close()
				body, err := io.ReadAll(response.Body)
				json.Unmarshal(body, &receivedResponse)
	*/
}

func JoinEventHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	glog.V(2).Info("JoinEventHandler: ", r)
	var response JoinEventResponse

	clientId := uuid.NewString()

	response.ClientId = clientId

	//empty socket, socket will be initialized in socket handlers
	room.LiveQuestionHandler.ClientsMapping[clientId] = nil
	room.LiveResultsHandler.ClientsMapping[clientId] = nil

	//TODO: Send question depending on the state.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ClientsResponseHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	glog.V(2).Info("ClientsResponseHandler: ", r)
	var response LiveResultsResponse

	// Pre-processing of the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		glog.Error("IO Request Body read failed", err)
		response.Error = "IO Request Body read failed" + err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	resultsCountMap, err := room.CollectClientResponse(bodyBytes)
	if err != nil {
		response.Error = err.Error()
	} else {
		response.LiveResults = resultsCountMap
	}

	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(response)
	responseBytes, err := json.Marshal(response)
	if err != nil {
		glog.Error(err)
	}
	w.Write(responseBytes)
}

func AddLiveQuestionHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	glog.V(2).Info("AddLiveQuestionHandler: ", r)
	// TODO: Handle the usage of roomId (API body contract) using r
	var response CreateQuestionResponse

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		glog.Error("IO Request Body read failed", err)
		response.Error = "IO Request Body read failed" + err.Error()
	} else {
		questionId, err := room.AddLiveQuestion(bodyBytes)
		response.QuestionId = questionId

		if err != nil {
			response.Error = err.Error()
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	//TODO: Check if data sent is correct??
	//TODO: check if the process is correct
	//Start sending go routine after 5 secs
	go func() {
		SendLiveQuestion(room, bodyBytes)
		room.SendLiveResponse(room.LiveResultsHandler)
	}()
}

//TODO: See if we need to send question multiple times??
func SendLiveQuestion(room *room.RoomInstance, question []byte) {

	//first register all the available clients
	room.LiveQuestionHandler.RegisterAllClients()

	room.LiveQuestionHandler.Broadcast <- question
}

func FetchCurrentStateHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	glog.V(2).Info("FetchCurrentStateHandler: ", r)
	var response FetchCurrentStateResponse

	response.State = room.FetchCurrentState()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func FetchLiveQuestionHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	glog.V(2).Info("FetchLiveQuestionHandler: ", r)
	var response FetchLiveQuestionResponse

	liveQuestion, err := room.FetchLiveQuestion()
	if err != nil {
		response.Error = err.Error()
	} else {
		response.Owner = liveQuestion.Owner
		response.Question = liveQuestion.Question
		response.Options = liveQuestion.Options
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func EndEventHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	glog.V(2).Info("EndEventHandler: ", r)

	//close the live results and live question handlers
	go func() {
		room.LiveResultsHandler.Close <- true
		room.LiveQuestionHandler.Close <- true
	}()

	var response EndEventResponse
	err := room.EndEvent()
	if err != nil {
		response.Error = err.Error()
	}
	w.WriteHeader((http.StatusOK))
	json.NewEncoder(w).Encode(response)
}

func MoveToNextQuestionHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	glog.V(2).Info("MoveToNextQuestionHandler: ", r)
	var response MoveToNextQuestionResponse
	err := room.MoveToNextQuestion()
	if err != nil {
		response.Error = err.Error()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	//TODO: Notify clients to navigate to next question

	//clear the registered clients map in LiveResultsHandler
	go room.LiveResultsHandler.UnRegisterAllClients()

	//close the response go routine
	go func() { room.StopSendingLiveResults <- true }()
}
