// This file consists of the handlers used by the Server
package rest

import (
	"interact/server/data"
	room "interact/server/room"
	"net/http"

	"encoding/json"

	"github.com/golang/glog"
	"github.com/google/uuid"

	// "github.com/gorilla/mux"
	"io"
)

func SetResponseMetadata(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, X-Auth-Token")
	w.Header().Set("Content-Type", "application/json")
}

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
	// Sample usecase to display text on webpage
	// vars := mux.Vars(r)
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div><div>%v</div>", "Interact",
	// 	"Application", vars)
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
	SetResponseMetadata(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ClientsResponseHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) bool {
	glog.V(2).Info("ClientsResponseHandler: ", r)
	var response LiveResultsResponse

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Pre-processing of the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		glog.Error("IO Request Body read failed", err)
		response.Error = "IO Request Body read failed" + err.Error()
		SetResponseMetadata(w)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return false
	}

	resultsCountMap, err := room.CollectClientResponse(bodyBytes)
	if err != nil {
		response.Error = err.Error()
		SetResponseMetadata(w)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return false
	} else {
		response.LiveResults = resultsCountMap
	}

	// json.NewEncoder(w).Encode(response)
	responseBytes, err := json.Marshal(response)
	if err != nil {
		glog.Error(err)
		response.Error = err.Error()
		SetResponseMetadata(w)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return false
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
	return true
}

func AddLiveQuestionHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	// TODO: Add Host Validation with host-id
	glog.V(2).Info("AddLiveQuestionHandler: ", r)
	var response CreateQuestionResponse

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		glog.Error("IO Request Body read failed", err)
		response.Error = "IO Request Body read failed" + err.Error()
		SetResponseMetadata(w)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	pollData, err := data.NewLivePollData(bodyBytes)
	if err != nil {
		response.Error = err.Error()
		SetResponseMetadata(w)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	questionId, err := room.AddLiveQuestion(pollData)
	if err != nil {
		response.Error = err.Error()
		SetResponseMetadata(w)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response.QuestionId = questionId

	// TODO: Add code to send the question to all clients using socket
	SetResponseMetadata(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	//TODO: check if the live response handler should be started instantly after sending the question
	bytesToSend, err := data.ConvertCurrQuestionToBytes(questionId, pollData)
	if err != nil {
		glog.Error("Question Conversion to Bytes failed ", err.Error())
		errResponse := ErrorResponse{
			Error: err.Error(),
		}
		bytesToSend, _ = json.Marshal(errResponse)
	}
	go func() {
		room.SendLiveQuestion(bytesToSend)
		room.SendLiveResponse(room.LiveResultsHandler)
	}()
}

func FetchCurrentStateHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	glog.V(2).Info("FetchCurrentStateHandler: ", r)
	var response FetchCurrentStateResponse

	response.State = room.FetchCurrentState()

	SetResponseMetadata(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func FetchLiveQuestionHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	glog.V(2).Info("FetchLiveQuestionHandler: ", r)
	var response FetchLiveQuestionResponse

	SetResponseMetadata(w)
	w.WriteHeader(http.StatusOK)

	liveQuestion, err := room.FetchLiveQuestion()
	if err != nil {
		response.Error = err.Error()
	} else {
		response.Owner = liveQuestion.Owner
		response.Question = liveQuestion.Question
		response.Options = liveQuestion.Options
	}

	json.NewEncoder(w).Encode(response)
}

func EndEventHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	// TODO: Add Host Validation with host-id
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
	SetResponseMetadata(w)
	w.WriteHeader((http.StatusOK))
	json.NewEncoder(w).Encode(response)
}

func MoveToNextQuestionHandler(w http.ResponseWriter, r *http.Request, room *room.RoomInstance) {
	glog.V(2).Info("MoveToNextQuestionHandler: ", r)
	var response MoveToNextQuestionResponse
	// TODO: Add Host Validation with host-id
	err := room.MoveToNextQuestion()
	if err != nil {
		response.Error = err.Error()
	}

	SetResponseMetadata(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	go func() {
		bytes, err := NotifyStateChangeBytes(UI_STATE_LOADING)
		if err != nil {
			glog.Error("notifyStateChangeBytes error: ", err.Error())
		}
		//clear the registered clients map in LiveResultsHandler
		room.NotifyClientsForNextQuestion(bytes)
		room.LiveResultsHandler.UnRegisterAllClients()
	}()
	//close the response go routine
	go func() { room.StopSendingLiveResults <- true }()
}

func NotifyStateChangeBytes(state int) ([]byte, error) {
	resp := NotifyStateChangeResponse {
		State: state,
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		glog.Error(err.Error())
	}
	return bytes, err
}