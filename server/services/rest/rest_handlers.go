// This file consists of the handlers used by the Server
package rest

import (
	"fmt"
	"github.com/golang/glog"
	socketio "github.com/googollee/go-socket.io"
	data "interact/server/config/data"
	"net/http"
)

func CreateInstanceHandler(w http.ResponseWriter, r *http.Request, socket *socketio.Server, room *data.RoomInstance) error {
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
	return nil
}
