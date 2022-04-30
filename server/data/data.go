// This file consists of the structures to store data for different eventTypes.
package data

type questionType int

const (
	wordAnswer questionType = iota
	singleCorrect
	multiCorrect
)

type LivePollData struct {
	// id denotes the Question Id
	id    int
	owner *string
	questionType
	question *string
	// options, answer - Used only incase of MCQs
	options []*string
	// Use bitmask/sortedOption String to store the answer
	answer    *string
	responses []*clientResponse
	// TODO: Add lock to handle the concurrency of collecting responses
}

// Instead of ResponseData, could use interface for word/mcq response since
// string would a return type for getResponse func. In that case we can't use
// the same struct for unmarshalling the clients API data response
// Use utils.go for any conversions
type ResponseData struct {
	wordResponse *string
	mcqResponse  *string
}

type clientResponse struct {
	clientId int
	ResponseData
}

type LivePollResults struct {
	questionId     int
	answerCountMap map[string]int
}
