// This file consists of the data definition of the roomInstance
package room

import (
	data "interact/server/data"
	socket "interact/server/services/socket"
)

type eventType int
type state int

const (
	LivePolls eventType = iota
)

const (
	WAITING_ON_HOST_FOR_QUESTION state = iota
	WAITING_ON_CLIENTS_FOR_RESPONSES
)

type RoomInstance struct {
	// roomId specifies the id of the room. eg: suffix at url
	roomId string
	eventType
	// pollsData consists of all the questions data of the event
	pollsData    []*data.LivePollData
	currentState state
	// currentQuestion specifies the data of the current question
	// it is added to the pollsData after the question is Done ie., the host
	// moves to the next question OR stops collecting the responses for the
	// question OR ends the event
	currentQuestion   *data.LivePollData
	numOfParticipants int
	// TODO: Check if lock is required while updating the numOfParticipates var
	ClientHandler *socket.ClientHandler
}

const defaultId = "default-room-id"

func NewRoomInstance() *RoomInstance {
	room := &RoomInstance{
		roomId:        defaultId,
		ClientHandler: socket.NewClientHandler(),
	}

	return room
}
