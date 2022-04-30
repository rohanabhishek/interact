// This file consists of the handlers used by the Server
package rest

import (
	"fmt"
	room "interact/server/room"
	"net/http"

	"encoding/json"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
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
	// Sample usecase to display text on webpage
	vars := mux.Vars(r)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div><div>%v</div>", "Interact",
		"Application", vars)
	if roomInstanceResponse.Error != "" {
		json.NewEncoder(w).Encode(roomInstanceResponse)
		return
	}

	// Pre-processing of the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		glog.Error("IO Request Body read failed", err)
		roomInstanceResponse.Error = "IO Request Body read failed" + err.Error()
		json.NewEncoder(w).Encode(roomInstanceResponse)
		return
	}

	err = room.SetRoomConfig(bodyBytes)
	if err != nil {
		roomInstanceResponse.Error = err.Error()
	}

	// TODO: set the status of the APIs appropriately in case of errors
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roomInstanceResponse)
}

func JoinEventHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	glog.V(2).Info("JoinEventHandler: ", r)
	var response JoinEventResponse
	response.ClientId = room.GetNewClientId()
	response.Error = ""

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

	// TODO: Add this client to socket, who can view the results
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
	// TODO: Add code to send the question to all clients using socket
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
	glog.V(2).Info("MoveToNextQuestionHandler: ", r)
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
}
