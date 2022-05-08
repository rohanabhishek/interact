// This file consists of the API responses format
package rest

import (
	"interact/server/room"
)

type CreateInstanceResponse struct {
	RoomId string `json:"roomId"`
	Error  string `json:"error,omitempty"`
}

type CreateQuestionResponse struct {
	QuestionId int    `json:"questionId"`
	Error      string `json:"error,omitempty"`
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
	Owner *string `json:"owner"`
	// TODO: Add QuestionType by moving it to someother package
	Question *string   `json:"question"`
	Options  []*string `json:"options"`
	Error    string    `json:"error,omitempty"`
}

type MoveToNextQuestionResponse struct {
	Error string `json:"error,omitempty"`
}

type EndEventResponse struct {
	Error string `json:"error,omitempty"`
}
