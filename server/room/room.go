// This file consists of the data definition of the roomInstance
package room

import (
	"encoding/json"
	"errors"
	data "interact/server/data"
	socket "interact/server/services/socket"
	"sync"
	"time"

	"github.com/golang/glog"
)

type EventType int
type State int

const (
	LivePolls EventType = iota
)

const (
	WAITING_ON_HOST_FOR_QUESTION State = iota
	SENDING_QUESTION_TO_CLIENTS
	COLLECTING_CLIENT_RESPONSES
	EVENT_END
)

type LiveSocketResponse struct {
	State    `json:"state"`
	Response interface{} `json:"response"`
	Error    error       `json:"error,omitempty"`
}

type RoomInstance struct {
	// roomId specifies the id of the room. eg: suffix at url
	roomId string
	hostId string
	// hostName specifies the host of the room
	HostName  string `json:"hostName"`
	EventType `json:"eventType"`
	// pollsData consists of all the questions data of the event
	pollsData    []*data.LivePollData
	currentState State
	// currentQuestion specifies the data of the current question
	// it is added to the pollsData after the question is Done ie., the host
	// moves to the next question OR stops collecting the responses for the
	// question OR ends the event
	currentQuestion        *data.LivePollData
	numOfParticipants      int
	SocketHandler          *socket.ClientHandler
	StopSendingLiveResults chan bool
	// qMutex handles sync for questions data
	qMutex sync.RWMutex
	// pMutex handles sync for participants count
	pMutex sync.Mutex
}

func NewRoomInstance(roomId string, hostId string) *RoomInstance {
	room := &RoomInstance{
		roomId:                 roomId,
		hostId:                 hostId,
		currentState:           WAITING_ON_HOST_FOR_QUESTION,
		currentQuestion:        nil,
		numOfParticipants:      0,
		StopSendingLiveResults: make(chan bool),
		SocketHandler:          socket.NewClientHandler(),
	}

	return room
}

func ValidateRoomUnMarshal(bytes []byte) error {
	tempRoom := new(RoomInstance)
	return tempRoom.UnMarshal(bytes)
}

func (room *RoomInstance) UnMarshal(bytes []byte) error {
	rawStructData := &struct {
		HostName  string      `json:"hostName"`
		EventType interface{} `json:"eventType"`
	}{}
	err := json.Unmarshal(bytes, rawStructData)
	if err != nil {
		errMsg := "roomInstance: Unmarshal failed " + err.Error()
		glog.Error(errMsg)
		return errors.New(errMsg)
	}

	if rawStructData.HostName == "" {
		errMsg := "empty HostName not allowed"
		glog.Error(errMsg)
		return errors.New(errMsg)
	}

	room.HostName = rawStructData.HostName
	switch rawStructData.EventType {
	case "LivePolls":
		room.EventType = LivePolls
	default:
		err := "EventType not found"
		glog.Errorln(err)
		return errors.New(err)
	}
	return nil
}

func (room *RoomInstance) GetRoomId() string {
	return room.roomId
}

func (room *RoomInstance) GetHostId() string {
	return room.hostId
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
	room.qMutex.Lock()
	defer room.qMutex.Unlock()
	room.pollsData = append(room.pollsData, room.currentQuestion)
	room.currentQuestion = nil
	room.currentState = WAITING_ON_HOST_FOR_QUESTION
	return nil
}

func (room *RoomInstance) SetRoomStateToCollectResponses() error {
	room.qMutex.Lock()
	defer room.qMutex.Unlock()
	if room.currentState != SENDING_QUESTION_TO_CLIENTS {

		glog.Error("error setting room state to collecting client responses")

		return errors.New("Cannot set state to collecting responses as current state is not sending question")
	}

	//set the current state
	room.currentState = COLLECTING_CLIENT_RESPONSES

	return nil
}

func (room *RoomInstance) AddLiveQuestion(pollData *data.LivePollData) (int, error) {
	room.qMutex.Lock()
	defer room.qMutex.Unlock()
	if room.currentState != WAITING_ON_HOST_FOR_QUESTION {
		glog.Error("AddLiveQuestion: Currenstate is not WAITING_ON_HOST_FOR_QUESTION")
		return -1, errors.New("current State is not WAITING_ON_HOST_FOR_QUESTION")
	}

	room.currentQuestion = pollData
	questionId := len(room.pollsData) + 1
	// start accepting the client's responses
	room.currentState = SENDING_QUESTION_TO_CLIENTS
	return questionId, nil
}

func (room *RoomInstance) GetNumOfParticipants() int {
	room.pMutex.Lock()
	defer room.pMutex.Unlock()
	return room.numOfParticipants
}

func (room *RoomInstance) GetNewClientId() int {
	room.pMutex.Lock()
	defer room.pMutex.Unlock()
	room.numOfParticipants += 1
	return room.numOfParticipants
}

func (room *RoomInstance) CollectClientResponse(args []byte) (map[string]int, error) {
	// TODO: We can use a channel to take input the decoded object from the API
	// handler and process the object directly
	room.qMutex.RLock()
	defer room.qMutex.RUnlock()
	if room.currentState != COLLECTING_CLIENT_RESPONSES {
		glog.Error("collectClientResponse: Currenstate is not WAITING_ON_CLIENTS_FOR_RESPONSES")
		return nil, errors.New("current State is not WAITING_ON_CLIENTS_FOR_RESPONSES")
	}

	return room.currentQuestion.CollectClientResponse(args)
}

func (room *RoomInstance) FetchCurrentState() State {
	room.qMutex.RLock()
	defer room.qMutex.RUnlock()
	return room.currentState
}

func (room *RoomInstance) FetchLiveQuestion() (*data.LivePollData, error) {
	room.qMutex.RLock()
	defer room.qMutex.RUnlock()
	if room.currentState != COLLECTING_CLIENT_RESPONSES {
		glog.Error("collectClientResponse: Currenstate is not WAITING_ON_CLIENTS_FOR_RESPONSES")
		return nil, errors.New("current State is not WAITING_ON_CLIENTS_FOR_RESPONSES")
	}

	return &data.LivePollData{
		Owner:        room.currentQuestion.Owner,
		QuestionType: room.currentQuestion.QuestionType,
		Question:     room.currentQuestion.Question,
		Options:      room.currentQuestion.Options,
	}, nil
}

func (room *RoomInstance) EndEvent() error {
	//TODO: store it in database and make room nil
	room.qMutex.Lock()
	defer room.qMutex.Unlock()
	room.pollsData = append(room.pollsData, room.currentQuestion)
	room.currentQuestion = nil
	room.currentState = EVENT_END
	return nil
}

/*
- Handling the clients' responses which are processed after the Host triggers
MoveToNextQuestion is being done,
	- We have the state checks whenever we collect the responses. Since, we
		have handled the question's Data and currentState using the same RWLock
		the data consistency will be held b/w these two, thereby we would just send
		and error message as response for the API call, instead be
		preventing the server to crash because of any null pointers.
	- TODO(Rohan) - Handle the case where one question's response gets aggregated
		in other question.
*/

//TODO: See if we need to send question multiple times??
func (room *RoomInstance) SendLiveQuestion(question []byte) {
	glog.Info("sending live question to clients")

	room.SocketHandler.RegisterAllClients()
	room.SocketHandler.Broadcast <- question

	glog.Info("current live question sent")
}

// Send the state change while MoveToNextQuestion triggered by Host
func (room *RoomInstance) NotifyClientsForNextQuestion(state []byte) {
	glog.Info("sending state change to all clients")

	room.SocketHandler.RegisterAllClients()
	room.SocketHandler.Broadcast <- state

	glog.Info("state change update sent")
}

//function to write live responses, it broadcasts message every one sec
func (room *RoomInstance) SendLiveResponse(ch *socket.ClientHandler) {
	ticker := time.NewTicker(1 * time.Second)

	glog.Info("Started sending live responses in socket handler")
	//Send based on the revision, dont send duplicate data.
	resLength := 0
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			//get the live data
			if room.currentState != COLLECTING_CLIENT_RESPONSES {
				//We moved to the next question
				glog.Info("Stopped sending live results")
				return
			}
			data, len := room.currentQuestion.GetResponseStats()

			//only send if current length is greater
			if len > resLength {

				glog.Info("Broadcasting live data")

				responseBytes := room.GetSocketResponse(data, nil)

				ch.Broadcast <- responseBytes

				resLength = len
			} else if len < resLength {
				glog.Error("How responses length decreased??")
			}

		case <-room.StopSendingLiveResults:
			glog.Info("stopped sending live results and moving to next question")
			return
		}
	}
}

func (room *RoomInstance) GetSocketResponse(res interface{}, err error) []byte {
	response := LiveSocketResponse{
		State:    room.FetchCurrentState(),
		Response: res,
		Error:    err,
	}

	bytesToSend, err := json.Marshal(response)

	if err != nil {
		glog.Error("Response conversion to bytes failed ", err.Error())

		response = LiveSocketResponse{
			Error: err,
		}

		bytesToSend, _ = json.Marshal(response)
	}

	return bytesToSend
}
