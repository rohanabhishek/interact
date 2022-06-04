package data

import (
	"errors"

	"github.com/golang/glog"
)

func GetAnswerAsString(response ResponseData,
	questionType QuestionType) (string, error) {
	switch questionType {
	case WordAnswer:
		return response.GetWordResponse()
	// TODO: Add the handling of multi correct questions
	// based on the API response received
	case SingleCorrect:
		return response.GetMcqResponse()
	default:
		return "", errors.New("unknown questionType found")
	}
}

type LiveQuestionData struct {
	Owner        *string   `json:"owner"`
	QuestionId   int       `json:"questionId"`
	QuestionType string    `json:"questionType"`
	Question     *string   `json:"question"`
	Options      []*Option `json:"options"`
	Error        string    `json:"error,omitempty"`
}

// ConvertCurrQuestionToBytes converts the current question to bytes so as to
// send to clients on socket
func GetCurrentQuestion(questionId int, pollData *LivePollData) (*LiveQuestionData, error) {
	response := &LiveQuestionData{
		Owner:    pollData.Owner,
		Question: pollData.Question,
		Options:  pollData.Options,
	}
	switch pollData.QuestionType {
	case SingleCorrect:
		response.QuestionType = SingleCorrectQuesTypeStr
	default:
		err := "pollData QuestionType not found"
		glog.Error(err)
		return nil, errors.New(err)
	}
	response.QuestionId = questionId

	return response, nil
}
