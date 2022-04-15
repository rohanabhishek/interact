package rest

type CreateInstanceResponse struct {
	RoomId string `json:"roomId"`
	Error  string `json:"error"`
}

type CreateQuestionResponse struct {
	QuestionId int    `json:"questionId"`
	Error      string `json:"error"`
}
