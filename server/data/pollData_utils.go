package data

import (
	"errors"
)

//use this to close any go routines
type CloseChannel chan bool

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
