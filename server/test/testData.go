package test

import (
	"interact/server/room"
	"interact/server/services/rest"
)

var (
	DefaultCreateInstanceParams = `
	{
		"hostName": "IAmTheHost",
		"eventType": "LivePolls"
	}
	`

	expectedCreateInstanceResponseVal = rest.CreateInstanceResponse{
		RoomId: "default-room-id",
	}

	currStateWaitingOnHost = rest.FetchCurrentStateResponse{
		State: room.WAITING_ON_HOST_FOR_QUESTION,
	}

	currStateWaitingOnClients = rest.FetchCurrentStateResponse{
		State: room.WAITING_ON_CLIENTS_FOR_RESPONSES,
	}

	endEventState = rest.FetchCurrentStateResponse{
		State: room.EVENT_END,
	}

	SampleAddQuestionParams = `
	{
		"owner": "IAmTheHost",
		"questionType": "Single MCQ",
		"question": "What is the current mood?",
		"options": [
			"Happy",
			"Sad",
			"Neutral"
		]
	}
	`

	addQuestionResponseVal = rest.CreateQuestionResponse{
		QuestionId: 1,
	}

	SampleClientResponse = `
	{
		"clientId": 1,
		"response": "A"
	}
	`

	FmtSampleClientResponse = `
	{
		"clientId": %d,
		"response": "%v"
	}
	`

	SampleAddQuestion2Params = `
	{
		"owner": "IAmTheHost",
		"questionType": "Single MCQ",
		"question": "How's the Josh?",
		"options": [
			"Very High",
			"Netural",
			"Low"
		]
	}
	`
	addQuestion2ResponseVal = rest.CreateQuestionResponse{
		QuestionId: 2,
	}
)
