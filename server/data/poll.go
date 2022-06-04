// This file consists of the structures to store data for different eventTypes.
package data

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"github.com/golang/glog"
)

type QuestionType int

type Option struct {
	Idx    int    `json:"idx"`
	Option string `json:"option"`
}

const (
	WordAnswer QuestionType = iota
	SingleCorrect
	MultiCorrect
)

const (
	SingleCorrectQuesTypeStr string = "Single MCQ"
)

type LivePollData struct {
	// id denotes the Question Id
	id           int
	Owner        *string `json:"owner"`
	QuestionType `json:"questionType"`
	Question     *string `json:"question"`
	// options, answer - Used only incase of MCQs
	Options []*Option `json:"options"`
	// Use bitmask/sortedOption String to store the answer
	Answer          *string `json:"answer"`
	responses       []*ClientResponse
	resultsCountMap map[string]int
	mutex           sync.RWMutex
}

// Instead of ResponseData, could use interface for word/mcq response since
// string would a return type for getResponse func. In that case we can't use
// the same struct for unmarshalling the clients API data response
// Use utils.go for any conversions
type ResponseData struct {
	WordResponse *string `json:"wordResponse"`
	McqResponse  *string `json:"mcqResponse"`
}

type ClientResponse struct {
	ClientId     string `json:"clientId"`
	ResponseData `json:"responseData"`
}

type LiveResult struct {
	Count      int `json:"count"`
	Percentage int `json:"percentage"`
	Idx        int `json:"idx"`
	Answer     int `json:"answer,omitempty"`
}

func (clientResponse *ClientResponse) UnMarshal(bytes []byte, qtype QuestionType) error {
	rawStructData := &struct {
		ClientId string `json:"clientId"`
		// Accepting the response as string so that we can handle the wordAnswers as well.
		// Use the following syntax for mcq answers as A/ABC/1/123
		Response string `json:"response"`
	}{}
	err := json.Unmarshal(bytes, rawStructData)
	//glog.Infof("clientResponse: bytes %v, rawStructData %v", bytes, rawStructData)
	if err != nil {
		glog.Error("clientResponse: Unmarshal failed", err.Error())
		return err
	}

	clientResponse.ClientId = rawStructData.ClientId
	switch qtype {
	case SingleCorrect:
		var responseData ResponseData
		responseData.McqResponse = new(string)
		mcqResponse := rawStructData.Response
		responseData.McqResponse = &mcqResponse
		clientResponse.ResponseData = responseData
	}

	//glog.Infof("clientResponse Unmarshal: rawStructData: %v, clientResponse: %v", rawStructData, clientResponse)
	return nil
}

// GetWordResponse returns the WordResponse
func (data *ResponseData) GetWordResponse() (string, error) {
	if data.WordResponse == nil {
		return "", errors.New("word Response is nil")
	}
	return *data.WordResponse, nil
}

// GetWordResponse returns the McqResponse
func (data *ResponseData) GetMcqResponse() (string, error) {
	if data.McqResponse == nil {
		return "", errors.New("mcq Response is nil")
	}
	return *data.McqResponse, nil
}

func NewLivePollData(args []byte) (*LivePollData, error) {
	livePollData := new(LivePollData)
	livePollData.resultsCountMap = make(map[string]int)
	if err := livePollData.UnMarshal(args); err != nil {
		glog.Error("Unmarshal of LivePollData failed", err.Error())
		return nil, err
	}

	if livePollData.QuestionType == SingleCorrect {
		//initialize map with the options as key
		for _, element := range livePollData.Options {
			livePollData.resultsCountMap[strconv.Itoa(element.Idx)] = 0
		}

		for k, v := range livePollData.resultsCountMap {
			glog.Info(k, "value is", v)
		}
	}

	return livePollData, nil
}

func (pollData *LivePollData) UnMarshal(bytes []byte) error {
	rawStructData := &struct {
		Owner        *string     `json:"owner"`
		QuestionType interface{} `json:"questionType"`
		Question     *string     `json:"question"`
		// options, answer - Used only incase of MCQs
		Options []*Option `json:"options"`
		// Use bitmask/sortedOption String to store the answer
		Answer *string `json:"answer"`
	}{}
	err := json.Unmarshal(bytes, rawStructData)
	if err != nil {
		glog.Error("RoomInstance: Unmarshal failed", err.Error())
		return err
	}
	pollData.Owner = rawStructData.Owner
	pollData.Question = rawStructData.Question

	if rawStructData.Options != nil {
		pollData.Options = rawStructData.Options
	}

	if rawStructData.Answer != nil {
		pollData.Answer = rawStructData.Answer
	}

	switch rawStructData.QuestionType {
	case SingleCorrectQuesTypeStr:
		pollData.QuestionType = SingleCorrect
	}
	return nil
}

func (pollData *LivePollData) CollectClientResponse(apiResponse []byte) (map[string]int, error) {
	response := new(ClientResponse)
	if err := response.UnMarshal(apiResponse, pollData.QuestionType); err != nil {
		glog.Error("Unmarshal of ClientResponse failed", err.Error())
		return nil, err
	}
	pollData.mutex.Lock()
	defer pollData.mutex.Unlock()
	pollData.responses = append(pollData.responses, response)
	answer, err := GetAnswerAsString(response.ResponseData, pollData.QuestionType)
	if err != nil {
		glog.Errorln("collectClientResponse failed", err.Error())
		return nil, err
	}
	if count, exists := pollData.resultsCountMap[answer]; exists {
		pollData.resultsCountMap[answer] = count + 1
	} else {
		pollData.resultsCountMap[answer] = 1
	}
	glog.V(2).Infof("resultsCountMap %v", pollData.resultsCountMap)
	return pollData.resultsCountMap, nil
}

//TODO: both api and socket should use same function
func (pollData *LivePollData) GetLiveSocketResponse() ([]*LiveResult, int, error) {
	// TODO: Use utils.go and convert the response as per the UI's
	// frontend handler requirement which will be sent through the socket IO
	// We might need to send the Answer also, so as to display on UI
	pollData.mutex.RLock()
	defer pollData.mutex.RUnlock()
	liveResults, err := GetLiveResponse(pollData.resultsCountMap)
	return liveResults, len(pollData.responses), err
}

func GetLiveResponse(resultsCountMap map[string]int) ([]*LiveResult, error) {
	//TODO: take question type as variable and send idx or answer accordingly

	totalCount := 0

	var responses []*LiveResult

	for idx, count := range resultsCountMap {
		index, err := strconv.Atoi(idx)

		if err != nil {
			glog.Error("How index is not parsable", err)
			return nil, err
		}

		response := &LiveResult{
			Count:      count,
			Idx:        index,
			Percentage: 0,
		}

		totalCount += count

		responses = append(responses, response)
	}

	if totalCount == 0 {
		return responses, nil
	}

	for _, response := range responses {

		response.Percentage = (response.Count * 100 / totalCount)
	}

	return responses, nil
}
