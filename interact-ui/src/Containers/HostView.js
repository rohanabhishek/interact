import React, {
  useState,
  useEffect,
  useRef,
  useCallback,
  useContext,
} from "react";
import { UserContext } from "../UserContext.js";
import QuestionCard from "../Components/AddQuestion/Question";
import HostLiveResultsView from "./HostLiveResultsView";

const State = {
  loading: 0,
  question: 1,
  liveResults: 2,
  error: 3,
};

const intialRetryTime = 1000
const maxRetryTimeOut = 5000

const HostView = () => {
  const { contextDetails } = useContext(UserContext);
  let roomId = contextDetails.roomId;
  let clientId = contextDetails.userId;

  const ws = useRef(null);

  const [state, setState] = useState(State.question);
  const [pollData, setPollData] = useState(null);
  const [question, setQuestion] = useState(null);
  const [error, setError] = useState(null)

  //TODO: try to use a single function (make it as a wrapper)
  const connect = (ws, roomId, clientId, path, retryTime) =>{

    ws.current = new WebSocket(`ws://localhost:8080/${roomId}/${path}/${clientId}`);

    console.log(ws)

    let isConnected = false

    ws.current.onopen = (e) => {
      console.log("[open] Connection established" , e)
      isConnected = true
    }
 
    ws.current.onclose = (e) => {

      let currentRetryTime = retryTime

      if(isConnected){
        currentRetryTime = intialRetryTime
      }

      let nextRetryTime = currentRetryTime + Math.floor(Math.random() * intialRetryTime);

      console.log(nextRetryTime)

      if(nextRetryTime <= maxRetryTimeOut){
        console.log(`Socket is closed. Reconnect will be attempted in ${nextRetryTime} second.`, e.reason);
        setTimeout(() =>{
          connect(ws,roomId, clientId, path, 2*currentRetryTime)
        }, nextRetryTime)
      }
      else{
        console.error(`Socket is closed. Maximum retries are reached`, e.reason);
        setState(State.error)
        setError(e.reason)
      } 
   }
 
   ws.current.onerror = function(err) {
     console.error('Socket encountered error: ', err.message, 'Closing socket');
     ws.current.close();
   }
  }

  const setStateQuestion = useCallback((data) => {
    console.log(data);
    setQuestion(data);
    setState(State.liveResults);
  }, []);

  const changeStateToQuestion = useCallback(() => {
    setState(State.question);
  }, []);

  useEffect(() => {
    connect(ws, roomId, clientId, 'ws', intialRetryTime)
  }, []);

  useEffect(()=>{
    if(!ws.current) return

    ws.current.onmessage = (e) =>{
      console.log("event", e)

      const message = JSON.parse(e.data)

      console.log("message", message)

      if(message.error != null){
        setError(message.error)
        setState(State.error)
      }
      else if(message.state!=null){
        switch (message.state){
          //TODO: handle other cases
          case 0:
            setPollData(message.response)
            setState(State.liveResults)
            break
          case 2:
            setPollData(message.response)
            setState(State.liveResults)
            break
        }
      }
    }

  },[])

  return [
    state === State.question && <QuestionCard setState={setStateQuestion} />,
    state === State.liveResults && (
      <HostLiveResultsView
        key={3}
        question={question.question}
        options={question.options}
        count={pollData}
        loading={state === 0}
        roomId={roomId}
        changeStateToQuestion={changeStateToQuestion}
      />
    ),
  ];
};

export default HostView;
