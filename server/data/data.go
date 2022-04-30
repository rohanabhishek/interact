// This file consists of the structures to store data for different eventTypes.
package data

import (
	"encoding/json"
	"errors"
	"github.com/golang/glog"
	"sync"
)

type QuestionType int

const (
	WordAnswer QuestionType = iota
	SingleCorrect
	MultiCorrect
)

type LivePollData struct {
	// id denotes the Question Id
	id           int
	Owner        *string `json:"owner"`
	QuestionType `json:"questionType"`
	Question     *string `json:"question"`
	// options, answer - Used only incase of MCQs
	Options []*string `json:"options"`
	// Use bitmask/sortedOption String to store the answer
	Answer          *string `json:"answer"`
	responses       []*clientResponse
	resultsCountMap map[string]int
	mutex           sync.RWMutex
}

// Instead of ResponseData, could use interface for word/mcq response since
// string would a return type for getResponse func. In that case we can't use
// the same struct for unmarshalling the clients API data response
// Use utils.go for any conversions
type ResponseData struct {
	WordResponse *string `json:wordResponse`
	McqResponse  *string `json:mcqResponse`
}

type clientResponse struct {
	ClientId     int `json:clientId`
	ResponseData `json:"responseData"`
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
	if err := json.Unmarshal(args, livePollData); err != nil {
		glog.Error("Unmarshal of LivePollData failed", err.Error())
		return nil, err
	}
	return livePollData, nil
}

func (pollData *LivePollData) CollectClientResponse(apiResponse []byte) (map[string]int, error) {
	response := new(clientResponse)
	if err := json.Unmarshal(apiResponse, response); err != nil {
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
	return pollData.resultsCountMap, nil
}

func (pollData *LivePollData) getResponseStats() map[string]int {
	// TODO: Use utils.go and convert the response as per the UI's
	// frontend handler requirement which will be sent through the socket IO
	// We might need to send the Answer also, so as to display on UI
	pollData.mutex.RLock()
	defer pollData.mutex.RUnlock()
	return pollData.resultsCountMap
}
