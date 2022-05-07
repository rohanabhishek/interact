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
	WAITING_ON_CLIENTS_FOR_RESPONSES
	EVENT_END
)

type RoomInstance struct {
	// roomId specifies the id of the room. eg: suffix at url
	roomId string
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
	LiveResultsHandler     *socket.ClientHandler
	LiveQuestionHandler    *socket.ClientHandler
	StopSendingLiveResults chan bool
	// qMutex handles sync for questions data
	qMutex sync.RWMutex
	// pMutex handles sync for participants count
	pMutex sync.Mutex
}

func NewRoomInstance(roomId string) *RoomInstance {
	room := &RoomInstance{
		roomId:                 roomId,
		currentState:           WAITING_ON_HOST_FOR_QUESTION,
		currentQuestion:        nil,
		numOfParticipants:      0,
		LiveResultsHandler:     socket.NewClientHandler(),
		LiveQuestionHandler:    socket.NewClientHandler(),
		StopSendingLiveResults: make(chan bool),
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
	room.currentState = WAITING_ON_CLIENTS_FOR_RESPONSES
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
	if room.currentState != WAITING_ON_CLIENTS_FOR_RESPONSES {
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
	if room.currentState != WAITING_ON_CLIENTS_FOR_RESPONSES {
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

	//first register all the available clients
	room.LiveQuestionHandler.RegisterAllClients()

	room.LiveQuestionHandler.Broadcast <- question
}

// Send the state change while MoveToNextQuestion triggered by Host
func (room *RoomInstance) NotifyClientsForNextQuestion(state []byte) {
	// Send this to both handles since clients could be on either of them
	room.LiveQuestionHandler.RegisterAllClients()
	room.LiveQuestionHandler.Broadcast <- state

	room.LiveResultsHandler.Broadcast <- state
}

//function to write live responses, it broadcasts message every one sec
func (room *RoomInstance) SendLiveResponse(ch *socket.ClientHandler) {
	ticker := time.NewTicker(1 * time.Second)

	//Send based on the revision, dont send duplicate data.
	resLength := 0
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			//get the live data
			if room.currentState != WAITING_ON_CLIENTS_FOR_RESPONSES {
				//We moved to the next question
				glog.Info("Stopped sending live results")
				return
			}
			data, len := room.currentQuestion.GetResponseStats()

			//only send if current length is greater
			if len > resLength {
				responseBytes, err := json.Marshal(data)

				if err != nil {
					glog.Error(err)
				}

				ch.Broadcast <- responseBytes

				resLength = len
			} else if len < resLength {
				glog.Error("How responses length decreased??")
			}

		case <-room.StopSendingLiveResults:
			return
		}
	}
}
