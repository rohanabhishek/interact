package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"interact/server/data"
	"interact/server/services/rest"
	"interact/server/services/web"
	"io"
	"net/http"
	"reflect"
	"testing"
)

var (
	serverIP = "127.0.0.1:8080"
)

type tempLiveResultsResponse struct {
	LiveResults map[string]interface{} `json:"liveResults"`
	Error       string                 `json:"error,omitempty"`
}

func setupTest(t *testing.T) {
	webServer := web.NewWebServer(serverIP)
	go webServer.Run()
}

func UnMarshalResponse(t *testing.T, body io.ReadCloser, dataStruct interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&dataStruct)
	if err == io.EOF {
		t.Log("EOF error in UnMarshal Response for ", err.Error())
		return nil
	} else if err != nil {
		t.Errorf("decoder.Decode :%v", err.Error())
	}
	return err
}

/*
Sample way to test a API response
func TestCreateInstance(t *testing.T) {
	testTables := []struct {
		description                    string
		inputCreateInstanceParams      string
		expectedCreateInstanceResponse rest.CreateInstanceResponse
	}{
		{
			description:                    "Default CreateInstance API call",
			inputCreateInstanceParams:      DefaultCreateInstanceParams,
			expectedCreateInstanceResponse: expectedCreateInstanceResponseVal,
		},
	}

	for _, tc := range testTables {
		setupTest(t)
		t.Run(tc.description, func(t *testing.T) {
			reqBody := bytes.NewBuffer([]byte(tc.inputCreateInstanceParams))
			t.Logf("reqBody %v", reqBody)
			response, err := http.Post("http://localhost:8080/createEvent", "text/plain", reqBody)
			if err != nil {
				t.Fatalf("Error in POST Request : %v", err.Error())
			}

			// decoder := json.NewDecoder(response.Body)
			// var receivedResponse rest.CreateInstanceResponse
			// err = decoder.Decode(&receivedResponse)
			// if err != nil {
				// 	t.Errorf("decoder.Decode :%v", err.Error())
			// }
			var receivedResponse rest.CreateInstanceResponse
			UnMarshalResponse(t, response.Body, &receivedResponse)

			if receivedResponse != tc.expectedCreateInstanceResponse {
				t.Fatalf("Recevied response %v, didn't match with the expected response %v",
					receivedResponse, tc.expectedCreateInstanceResponse)
			}
		})
	}
}
*/

func TestSampleRunOfEvent(t *testing.T) {
	setupTest(t)
	t.Run("TestSampleRunOfEvent", func(t *testing.T) {
		// Host Calls CreateEvent
		inputCreateInstanceParams := DefaultCreateInstanceParams
		reqBody := bytes.NewBuffer([]byte(inputCreateInstanceParams))
		createEventURL := fmt.Sprintf("http://%v/%v", serverIP, "createEvent")
		response, err := http.Post(createEventURL, "text/plain", reqBody)
		if err != nil {
			t.Fatalf("Error in POST Request : %v", err.Error())
		}

		var receivedResponse rest.CreateInstanceResponse
		UnMarshalResponse(t, response.Body, &receivedResponse)

		if receivedResponse != expectedCreateInstanceResponseVal {
			t.Fatalf("Recevied response %v, didn't match with the expected response %v",
				receivedResponse, expectedCreateInstanceResponseVal)
		} else {
			t.Logf("CreateEvent API response matched with expected %v", expectedCreateInstanceResponseVal)
		}

		// Client 1, 2, 3 Join the Event
		joinEventURL := fmt.Sprintf("http://%v/default-room-id/%v", serverIP, "joinEvent")
		for index := 1; index < 4; index++ {
			response, err = http.Post(joinEventURL, "text/plain", nil)
			if err != nil {
				t.Fatalf("Error in POST Request : %v", err.Error())
			}

			var joinEventResponse rest.JoinEventResponse
			UnMarshalResponse(t, response.Body, &joinEventResponse)
			// expectedJoinEventResponse := rest.JoinEventResponse{
			// 	ClientId: index,
			// }

			// if joinEventResponse != expectedJoinEventResponse {
			// 	t.Fatalf("Recevied response %v, didn't match with the expected response %v",
			// 		joinEventResponse, expectedJoinEventResponse)
			// } else {
			// 	t.Logf("JoinEvent API response matched with expected %v", expectedJoinEventResponse)
			// }
		}

		// Clients call fetchCurrentState and expect WAITING_ON_HOST_FOR_QUESTION as response
		fetchCurrentStateURL := fmt.Sprintf("http://%v/default-room-id/%v", serverIP, "fetchCurrentState")
		for index := 1; index < 4; index++ {
			response, err = http.Get(fetchCurrentStateURL)
			if err != nil {
				t.Fatalf("Error in Get Request : %v", err.Error())
			}

			var fetchCurrStateResponse rest.FetchCurrentStateResponse
			err = UnMarshalResponse(t, response.Body, &fetchCurrStateResponse)
			if err != nil {
				t.Fatalf("Unmarshal for fetchCurrentState failed %v", fetchCurrStateResponse)
			}
			expectedFetchCurrentStateResponse := currStateWaitingOnHost

			if fetchCurrStateResponse != expectedFetchCurrentStateResponse {
				t.Fatalf("Recevied response %v, didn't match with the expected response %v",
					fetchCurrStateResponse, expectedFetchCurrentStateResponse)
			} else {
				t.Logf("fetchCurrentState API response matched with expected %v", expectedFetchCurrentStateResponse)
			}
		}

		// Host adds the question
		addLiveQuestionURL := fmt.Sprintf("http://%v/default-room-id/%v", serverIP, "addLiveQuestion")
		liveQuestionRequestBody := bytes.NewBuffer([]byte(SampleAddQuestionParams))
		response, err = http.Post(addLiveQuestionURL, "text/plain", liveQuestionRequestBody)
		if err != nil {
			t.Fatalf("Error in POST Request : %v", err.Error())
		}

		var addQuestionResponse rest.CreateQuestionResponse
		err = UnMarshalResponse(t, response.Body, &addQuestionResponse)
		if err != nil {
			t.Fatalf("Unmarshal for createQuestion failed %v", addQuestionResponse)
		}
		expectedAddQuestionResponse := addQuestionResponseVal

		if addQuestionResponse != expectedAddQuestionResponse {
			t.Fatalf("Recevied response %v, didn't match with the expected response %v",
				addQuestionResponse, expectedAddQuestionResponse)
		} else {
			t.Logf("addLiveQuestion API response matched with expected %v", expectedAddQuestionResponse)
		}

		// Clients respond with answer
		sendResponseURL := fmt.Sprintf("http://%v/default-room-id/%v", serverIP, "sendResponse")
		clientResponses := make([]map[interface{}]interface{}, 3)
		clientResponses[0] =
			map[interface{}]interface{}{"clientId": 1, "response": "A"}
		clientResponses[1] =
			map[interface{}]interface{}{"clientId": 3, "response": "C"}
		clientResponses[2] =
			map[interface{}]interface{}{"clientId": 2, "response": "B"}
		for index, val := range clientResponses {
			clientResponseStr := fmt.Sprintf(FmtSampleClientResponse, val["clientId"], val["response"])
			// t.Logf("clientResponseStr %v", clientResponseStr)
			sendResponseBody := bytes.NewBuffer([]byte(clientResponseStr))
			response, err = http.Post(sendResponseURL, "text/plain", sendResponseBody)
			if err != nil {
				t.Fatalf("Error in POST Request : %v", err.Error())
			}
			var respClientResponse tempLiveResultsResponse
			// respClientResponse := make(map[string]interface{})

			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("body ReadAll err: %v", err.Error())
			}
			// else {
			// 	t.Logf("body %v", body)
			// }
			err = json.Unmarshal(body, &respClientResponse)
			if err != nil {
				t.Fatalf("Unmarshal for map[string]interface failed response %v, error: %v", respClientResponse, err.Error())
			} else {
				t.Logf("respClient %v", respClientResponse)
			}
			response.Body.Close()
			var clientRespFinal rest.LiveResultsResponse
			clientRespFinal.LiveResults = make(map[string]int)
			liveResultsMap := respClientResponse.LiveResults
			for key, val := range liveResultsMap {
				// t.Logf("key %v, val %v", key, val)
				clientRespFinal.LiveResults[key] = int(val.(float64))
			}

			expResp1 := map[string]int{
				"A": 1,
			}

			expResp2 := map[string]int{
				"A": 1,
				"C": 1,
			}

			expResp3 := map[string]int{
				"A": 1,
				"C": 1,
				"B": 1,
			}

			var expectedResponse map[string]int
			switch index {
			case 0:
				expectedResponse = expResp1
			case 1:
				expectedResponse = expResp2
			case 2:
				expectedResponse = expResp3
			}

			if !reflect.DeepEqual(expectedResponse, clientRespFinal.LiveResults) {
				t.Fatalf("Recevied response %v, didn't match with the expected response %v",
					clientRespFinal.LiveResults, expectedResponse)
			} else {
				t.Logf("sendResponse API response %v matched with expected %v", clientRespFinal.LiveResults, expectedResponse)
				// t.Log(clientRespFinal)
			}
		}

		// client:4 joins after responses collection started
		response, err = http.Post(joinEventURL, "text/plain", nil)
		if err != nil {
			t.Fatalf("Error in POST Request : %v", err.Error())
		}

		var joinEventResponse rest.JoinEventResponse
		UnMarshalResponse(t, response.Body, &joinEventResponse)
		// expectedJoinEventResponse := rest.JoinEventResponse{
		// 	ClientId: 4,
		// }

		// if joinEventResponse != expectedJoinEventResponse {
		// 	t.Fatalf("Recevied response %v, didn't match with the expected response %v",
		// 		joinEventResponse, expectedJoinEventResponse)
		// }

		//client:4 fetches the current state
		response, err = http.Get(fetchCurrentStateURL)
		if err != nil {
			t.Fatalf("Get Request error %v ", err.Error())
		}
		var fetchCurrStateResponse rest.FetchCurrentStateResponse
		err = UnMarshalResponse(t, response.Body, &fetchCurrStateResponse)
		if err != nil {
			t.Fatalf("Unmarshal for fetchCurrentState failed %v", fetchCurrStateResponse)
		}

		expectedFetchCurrentStateResponse := currStateWaitingOnClients

		if fetchCurrStateResponse != expectedFetchCurrentStateResponse {
			t.Fatalf("Recevied response %v, didn't match with the expected response %v",
				fetchCurrStateResponse, expectedFetchCurrentStateResponse)
		} else {
			t.Logf("fetchCurrentState API response matched with expected %v", expectedFetchCurrentStateResponse)
		}

		//client:4 fetches the live question
		fetchLiveQuestionURL := fmt.Sprintf("http://%v/default-room-id/%v", serverIP, "fetchLiveQuestion")
		response, err = http.Get(fetchLiveQuestionURL)
		if err != nil {
			t.Fatalf("Get Request error %v ", err.Error())
		}
		var fetchLiveQuestionResponse rest.FetchLiveQuestionResponse
		err = UnMarshalResponse(t, response.Body, &fetchLiveQuestionResponse)
		if err != nil {
			t.Fatalf("Unmarshal for fetchLiveQuestion failed %v", fetchLiveQuestionResponse)
		}

		poll := new(data.LivePollData)
		err = poll.UnMarshal([]byte(SampleAddQuestionParams))
		if err != nil {
			t.Fatalf("Unmarshal of Questiondata failed %v", err.Error())
		}

		compareOk := ((*(fetchLiveQuestionResponse.Owner) == *(poll.Owner)) && (*(fetchLiveQuestionResponse.Question) == *(poll.Question)) && compareArrPtr(fetchLiveQuestionResponse.Options, poll.Options))
		if !compareOk {
			t.Fatalf("Recevied response %v, didn't match with the expected response %v",
				fetchLiveQuestionResponse, expectedFetchCurrentStateResponse)
		} else {
			t.Logf("fetchLiveQuestion API response matched with expected %v", fetchLiveQuestionResponse)
		}

		// client:4 responds with answer
		clientResponseVal :=
			map[interface{}]interface{}{"clientId": 4, "response": "A"}

		clientResponseStr := fmt.Sprintf(FmtSampleClientResponse, clientResponseVal["clientId"], clientResponseVal["response"])
		sendResponseBody := bytes.NewBuffer([]byte(clientResponseStr))
		response, err = http.Post(sendResponseURL, "text/plain", sendResponseBody)
		if err != nil {
			t.Fatalf("Error in POST Request : %v", err.Error())
		}
		var respClientResponse tempLiveResultsResponse
		err = UnMarshalResponse(t, response.Body, &respClientResponse)
		if err != nil {
			t.Fatalf("Unmarshal for map[string]int failed %v", respClientResponse)
		}

		var clientRespFinal rest.LiveResultsResponse
		clientRespFinal.LiveResults = make(map[string]int)
		liveResultsMap := respClientResponse.LiveResults
		for key, val := range liveResultsMap {
			// t.Logf("key %v, val %v", key, val)
			clientRespFinal.LiveResults[key] = int(val.(float64))
		}

		expectedResponse := map[string]int{
			"A": 2,
			"C": 1,
			"B": 1,
		}

		if !reflect.DeepEqual(expectedResponse, clientRespFinal.LiveResults) {
			t.Fatalf("Recevied response %v, didn't match with the expected response %v",
				clientRespFinal, expectedResponse)
		} else {
			t.Logf("sendResponse API response %v matched with expected %v", clientRespFinal.LiveResults, expectedResponse)
		}

		// Host decides to move to next question
		moveToNextQuestionURL := fmt.Sprintf("http://%v/default-room-id/%v", serverIP, "nextLiveQuestion")
		response, err = http.Post(moveToNextQuestionURL, "text/plain", nil)
		if err != nil {
			t.Fatalf("Error in POST Request : %v", err.Error())
		}
		var moveToNextQuestionResponse rest.MoveToNextQuestionResponse
		err = UnMarshalResponse(t, response.Body, &moveToNextQuestionResponse)
		if err != nil {
			t.Fatalf("Unmarshal for fetchLiveQuestion failed %v", moveToNextQuestionResponse)
		}

		var moveToNextQuestionDefaultResponse rest.MoveToNextQuestionResponse

		if moveToNextQuestionResponse != moveToNextQuestionDefaultResponse {
			t.Fatalf("Recevied response %v, didn't match with the expected response %v",
				moveToNextQuestionResponse, moveToNextQuestionDefaultResponse)
		} else {
			t.Logf("moveToNextQuestion API response matched with expected %v", moveToNextQuestionResponse)
		}

		// Host adds a new question
		liveQuestionRequestBody = bytes.NewBuffer([]byte(SampleAddQuestion2Params))
		response, err = http.Post(addLiveQuestionURL, "text/plain", liveQuestionRequestBody)
		if err != nil {
			t.Fatalf("Error in POST Request : %v", err.Error())
		}

		var addQuestion2Response rest.CreateQuestionResponse
		err = UnMarshalResponse(t, response.Body, &addQuestion2Response)
		if err != nil {
			t.Fatalf("Unmarshal for createQuestion failed %v", addQuestion2Response)
		}
		expectedAddQuestionResponse = addQuestion2ResponseVal

		if addQuestion2Response != expectedAddQuestionResponse {
			t.Fatalf("Recevied response %v, didn't match with the expected response %v",
				addQuestion2Response, expectedAddQuestionResponse)
		} else {
			t.Logf("addLiveQuestion API response matched with expected %v", expectedAddQuestionResponse)
		}

		// skipping the client responses apart as the same collectresponses is tested
		// above

		// Host decides to end the event
		endEventURL := fmt.Sprintf("http://%v/default-room-id/%v", serverIP, "endEvent")
		response, err = http.Post(endEventURL, "text/plain", nil)
		if err != nil {
			t.Fatalf("Error in POST Request : %v", err.Error())
		}
		var endEventResponse rest.EndEventResponse
		err = UnMarshalResponse(t, response.Body, &endEventResponse)
		if err != nil {
			t.Fatalf("Unmarshal for fetchLiveQuestion failed %v", endEventResponse)
		}

		var endEventDefaultRespone rest.EndEventResponse

		if endEventResponse != endEventDefaultRespone {
			t.Fatalf("Recevied response %v, didn't match with the expected response %v",
				endEventResponse, endEventDefaultRespone)
		} else {
			t.Logf("endEvent API response matched with expected %v", endEventDefaultRespone)
		}

		// Fetch the current state of the event
		response, err = http.Get(fetchCurrentStateURL)
		if err != nil {
			t.Fatalf("Get request error %v ", err)
		}
		var fetchFinalStateResponse rest.FetchCurrentStateResponse
		err = UnMarshalResponse(t, response.Body, &fetchFinalStateResponse)
		if err != nil {
			t.Fatalf("Unmarshal for fetchCurrentState failed %v", fetchFinalStateResponse)
		}

		expectedFetchCurrentStateResponse = endEventState

		if fetchFinalStateResponse != expectedFetchCurrentStateResponse {
			t.Fatalf("Recevied response %v, didn't match with the expected response %v",
				fetchFinalStateResponse, expectedFetchCurrentStateResponse)
		} else {
			t.Logf("fetchCurrentState API response matched with expected %v", expectedFetchCurrentStateResponse)
		}
	})

}
