// This file consists of the API responses format
package rest

import (
	"interact/server/room"
)

const (
	UI_STATE_LOADING  = 0
	UI_STATE_QUESTION = 1
	UI_STATE_RESULTS  = 2
)

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type NotifyStateChangeResponse struct {
	State int `json:"stateChange"`
}

type CreateInstanceResponse struct {
	RoomId string `json:"roomId"`
	HostId string `json:"hostId"`
	Error  string `json:"error,omitempty"`
}

type CreateQuestionResponse struct {
	QuestionId int    `json:"questionId"`
	Error      string `json:"error,omitempty"`
}

type ClientsSendQuestionResponse struct {
	QuestionId   int     `json:"questionId"`
	QuestionType string  `json:"questionType"`
	Question     *string `json:"question"`
	// options, answer - Used only incase of MCQs
	Options []*string `json:"options"`
}

type JoinEventResponse struct {
	ClientId string `json:"clientId"`
	Error    string `json:"error,omitempty"`
}

type LiveResultsResponse struct {
	LiveResults map[string]int `json:"liveResults"`
	Error       string         `json:"error,omitempty"`
}

// TODO: Can merge FetchCurrentState into the JoinEventResponse
type FetchCurrentStateResponse struct {
	State room.State `json:"state"`
}

type FetchLiveQuestionResponse struct {
	Owner        *string   `json:"owner"`
	QuestionId   int       `json:"questionId"`
	QuestionType string    `json:"questionType"`
	Question     *string   `json:"question"`
	Options      []*string `json:"options"`
	Error        string    `json:"error,omitempty"`
}

type MoveToNextQuestionResponse struct {
	Error string `json:"error,omitempty"`
}

type EndEventResponse struct {
	Error string `json:"error,omitempty"`
}
