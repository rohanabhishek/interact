// This file consists of the handlers used by the Server
package rest

import (
	"encoding/json"
	"fmt"
	data "interact/server/config/data"
	"io"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

func CreateInstanceHandler(w http.ResponseWriter, r *http.Request, room *data.RoomInstance,
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
	if roomInstanceResponse.Error != "" {
		json.NewEncoder(w).Encode(roomInstanceResponse)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		glog.Error("IO Request Body read failed", err)
		roomInstanceResponse.Error = "IO Request Body read failed" + err.Error()
		json.NewEncoder(w).Encode(roomInstanceResponse)
		return
	}

	err = room.SetRoomConfig(bodyBytes)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	roomInstanceResponse.Error = errMsg
	json.NewEncoder(w).Encode(roomInstanceResponse)
}

func ClientsResponseHandler(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) {

}

func AddLiveQuestionHandler(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) {
	glog.V(2).Info("AddLiveQuestion: ", r)
	// TODO: Handle the usage of roomId (API body contract) using r
	// vars := mux.Vars(r)
	var response CreateQuestionResponse

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		glog.Error("IO Request Body read failed", err)
		response.Error = "IO Request Body read failed" + err.Error()
	} else {
		questionId, err := room.AddLiveQuestion(bodyBytes)
		response.QuestionId = questionId

		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		response.Error = errMsg
	}
	json.NewEncoder(w).Encode(response)
}

func FetchLiveQuestionHandler(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) {
	// TODO
	vars := mux.Vars(r)
	fmt.Fprintf(w, "<h1>%s</h1><div>%v</div>", "Interact", vars)
}

func EndEventHandler(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) {
	// TODO
	vars := mux.Vars(r)
	fmt.Fprintf(w, "<h1>%s</h1><div>%v</div>", "Interact", vars)
}

func MoveToNextQuestionHandler(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) {
	// TODO
	vars := mux.Vars(r)
	fmt.Fprintf(w, "<h1>%s</h1><div>%v</div>", "Interact", vars)
}
