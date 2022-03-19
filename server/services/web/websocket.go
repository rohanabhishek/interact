// This file consists of using Socket.IO package
package web

import (
	"github.com/golang/glog"
	socketio "github.com/googollee/go-socket.io"
)

// Creates a new server of Socket.IO
func NewWebSocket() *socketio.Server {
	return socketio.NewServer(nil)
}

// Event Handling of the Socket.IO Server
func EventHandling(ws *socketio.Server) {
	// TODO(Rohan): Use the feature of Rooms in Socket.IO while broadcasting
	// to all the clients in a Poll-Instance
	// ws.OnConnect("/", func(s socketio.Conn) error {
	// 	s.SetContext("")
	// 	glog.Info("connected:", s.ID())
	// 	return nil
	// })

	// ws.OnEvent("/client", "event", func(s socketio.Conn, msg string) {
	// 	glog.V(2).Info("[Event: /client] With message: event", msg)
	// 	s.Emit("reply", "have "+msg)
	// })

	// ws.OnDisconnect("/", func(s socketio.Conn, msg string) {
	// 	glog.Error("Connection Closed", msg)
	// })
}
