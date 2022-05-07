// This file is to run some basic tests and could be run directly
// go run tester.go
package main

import (
	// "bytes"
	"encoding/json"
	"flag"
	"fmt"
	"interact/server/data"
	"interact/server/room"
	"interact/server/services/rest"
	"interact/server/test"

	"github.com/golang/glog"
)

func main() {
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Parse()

	// roomInstance Unmarshal
	tempRoom := room.NewRoomInstance("default-id")
	reqBytes := []byte(test.DefaultCreateInstanceParams)
	glog.Infof("reqBytes %v", reqBytes)
	err := tempRoom.UnMarshal(reqBytes)
	if err != nil {
		glog.Error(err.Error())
	} else {
		glog.Infof("Unmarshal successful, Room: %v", tempRoom)
	}

	// LivePollData Unmarshal
	tempPoll, err := data.NewLivePollData([]byte(test.SampleAddQuestionParams))
	if err != nil {
		glog.Error("Unmarshal failed NewLivePollData", err.Error())
	} else {
		glog.Infof("Unmarshal successful, Poll: %v", tempPoll)
	}

	// ClientResponse Unmarshal
	clientResp := new(data.ClientResponse)
	err = clientResp.UnMarshal([]byte(test.SampleClientResponse), data.SingleCorrect)
	if err != nil {
		glog.Error("Unmarshal failed ClientResponse", err.Error())
	} else {
		glog.Infof("Unmarshal successful, ClientResponse: %v", clientResp)
	}

	// ClientResponse Unmarshal
	clientResp = new(data.ClientResponse)
	sampleClientResponse := fmt.Sprintf(test.FmtSampleClientResponse, 1, "A")
	err = clientResp.UnMarshal([]byte(sampleClientResponse), data.SingleCorrect)
	if err != nil {
		glog.Error("Unmarshal failed ClientResponse", err.Error())
	} else {
		mcqResponse, _ := clientResp.GetMcqResponse()
		glog.Infof("Unmarshal successful, FmtSampleClientResponse: %v, mcqResponse: %v", clientResp, mcqResponse)
	}

	liveResults := new(rest.LiveResultsResponse)
	liveResults.LiveResults = map[string]int{
		"A": 1,
	}
	resultsBytes, err := json.Marshal(liveResults)
	if err != nil {
		glog.Error("liveResults: Json marshal failed, resultsBytes: %v, err: %v", resultsBytes, err)
	}

	glog.Info("bytes: %v", resultsBytes)
	// b := new(bytes.Buffer)
	// err = json.NewEncoder(b).Encode(resultsBytes)
	// if err !=  nil {
	// 	glog.Error("Json NewEncoder error %v", err)
	// }

	// glog.Info("bytes: %v", resultsBytes)
	results := make(map[string]interface{})
	err = json.Unmarshal(resultsBytes, &results)
	if err != nil {
		glog.Error("liveResults: Json unmarshal failed, results: %v, err: %v", results, err)
	}
	glog.Info(results)
}
