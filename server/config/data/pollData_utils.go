package data

import (
	"errors"
)

func GetAnswerAsString(response ResponseData,
	questionType QuestionType) (string, error) {
	switch questionType {
	case WordAnswer:
		return response.GetWordResponse()
	// TODO: Add the handling of single, multi correction questions
	// based on the API response received
	default:
		return "", errors.New("unknown questionType found")
	}
}
