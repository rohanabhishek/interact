// This file consists of the API responses format
package rest

import "interact/server/room"

type CreateInstanceResponse struct {
	RoomId string `json:"roomId"`
	Error  string `json:"error"`
}

type CreateQuestionResponse struct {
	QuestionId int    `json:"questionId"`
	Error      string `json:"error"`
}

type JoinEventResponse struct {
	ClientId int    `json:"clientId"`
	Error    string `json:"error"`
}

type LiveResultsResponse struct {
	LiveResults map[string]int `json:"liveResults"`
	Error       string         `json:"error"`
}

type FetchCurrentStateResponse struct {
	State room.State `json:"state"`
}

type FetchLiveQuestionResponse struct {
	Owner *string `json:"owner"`
	// TODO: Add QuestionType by moving it to someother package
	Question *string   `json:"question"`
	Options  []*string `json:"options"`
	Error    string    `json:"error"`
}

type MoveToNextQuestionResponse struct {
	Error string `json:"error"`
}

type EndEventResponse struct {
	Error string `json:"error"`
}