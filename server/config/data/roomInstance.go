// This file consists of the data definition of the roomInstance
package data

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/golang/glog"
)

type EventType int
type state int

const (
	LivePolls EventType = iota
)

const (
	WAITING_ON_HOST_FOR_QUESTION state = iota
	WAITING_ON_CLIENTS_FOR_RESPONSES
)

type RoomInstance struct {
	// roomId specifies the id of the room. eg: suffix at url
	roomId string
	// hostName specifies the host of the room
	HostName  string `json:"hostName"`
	EventType `json:"eventType"`
	// pollsData consists of all the questions data of the event
	pollsData    []*LivePollData
	currentState state
	// currentQuestion specifies the data of the current question
	// it is added to the pollsData after the question is Done ie., the host
	// moves to the next question OR stops collecting the responses for the
	// question OR ends the event
	currentQuestion   *LivePollData
	numOfParticipants int
	mutex             sync.Mutex
}

func NewRoomInstance() *RoomInstance {
	// TODO: Add args eventType
	roomInstance := new(RoomInstance)
	roomInstance.roomId = "default-room-id"
	roomInstance.currentState = WAITING_ON_HOST_FOR_QUESTION
	roomInstance.currentQuestion = nil
	// Considering server separate from the clients
	roomInstance.numOfParticipants = 0
	return roomInstance
}

func (room *RoomInstance) SetRoomConfig(bytes []byte) error {
	tempRoom := new(RoomInstance)
	if err := json.Unmarshal(bytes, tempRoom); err != nil {
		glog.Error("RoomInstance: Unmarshal failed")
		return err
	}
	if tempRoom.HostName != "" {
		return room.SetHostName(tempRoom.HostName)
	}
	return nil
}

func (room *RoomInstance) GetRoomId() string {
	if room.roomId == "" {
		glog.Errorln("GetRoomId: Empty RoomId")
	}
	return room.roomId
}

func (room *RoomInstance) SetHostName(hostName string) error {
	if room.HostName != "" {
		errMsg := "SetHostName: HostName not empty. Attempt to Overwrite it"
		glog.Errorln(errMsg)
		return errors.New(errMsg)
	}
	room.HostName = hostName
	return nil
}

// At the start of the event, AddLiveQuestion is directly invoked
// After the event has started, to add a new question, MoveToNextQuestion,
// AddLiveQuestion are invoked
func (room *RoomInstance) MoveToNextQuestion() error {
	room.pollsData = append(room.pollsData, room.currentQuestion)
	room.currentQuestion = nil
	room.currentState = WAITING_ON_HOST_FOR_QUESTION
	return nil
}

func (room *RoomInstance) AddLiveQuestion(args []byte) (int, error) {
	if room.currentState != WAITING_ON_HOST_FOR_QUESTION {
		glog.Error("AddLiveQuestion: Currenstate is not WAITING_ON_HOST_FOR_QUESTION")
		return -1, errors.New("current state is not WAITING_ON_HOST_FOR_QUESTION")
	}

	questionId := len(room.pollsData) + 1
	var err error
	room.currentQuestion, err = NewLivePollData(args)
	if err != nil {
		return -1, err
	}
	return questionId, nil
}

func (room *RoomInstance) GetNumOfParticipants() int {
	return room.numOfParticipants
}

func (room *RoomInstance) GetNewClientId() int {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	room.numOfParticipants += 1
	return room.numOfParticipants
}
