// This file consists of the handlers used by the Server
package rest

import (
	"fmt"
	data "interact/server/config/data"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

func CreateInstanceHandler(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) {

	/*
		Usage of r:
		r.Method  // request method
		r.URL     // request URL
		r.Header  // request headers
		r.Body    // request body
		https://pkg.go.dev/net/http#Request
	*/
	glog.Info("CreateInstanceHandler: ", r)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "Interact", "Application")
}

func ClientsResponseHandler(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "<h1>%s</h1><div>%v</div>", "Interact", vars)
}

func AddLiveQuestion(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "<h1>%s</h1><div>%v</div>", "Interact", vars)
}

func FetchLiveQuestion(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "<h1>%s</h1><div>%v</div>", "Interact", vars)
}

func EndEvent(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "<h1>%s</h1><div>%v</div>", "Interact", vars)
}

func MoveToNextQuestion(w http.ResponseWriter, r *http.Request, room *data.RoomInstance) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "<h1>%s</h1><div>%v</div>", "Interact", vars)
}
