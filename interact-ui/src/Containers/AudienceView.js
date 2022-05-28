import React, {
  useState,
  useEffect,
  useRef,
  useCallback,
  useContext,
} from "react";
import AudienceLiveResultsView from "./AudienceLiveResultsView";
import AudienceQuestionView from "./AudienceQuestionView";
import App from "../App";
import { UserContext } from "../UserContext.js";

//TODO: Add poll end and question end state
const State = {
  loading: 0,
  question: 1,
  liveResults: 2,
  error: 3,
};

const intialRetryTime = 1000
const maxRetryTimeOut = 5000

const AudienceView = () => {
  const { contextDetails } = useContext(UserContext);
  let roomId = contextDetails.roomId;
  let clientId = contextDetails.userId;
  // let isHost = contextDetails.isHost;
  // let navigate = useNavigate();

  const ws = useRef(null);

  const [state, setState] = useState(State.loading);
  const [pollData, setPollData] = useState(null);
  const [question, setQuestion] = useState(null);
  const [error, setError] = useState(null)


  const connect = (ws, roomId, clientId, path, retryTime ) =>{

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

  const setStatePollData = useCallback((data) => {
    console.log(data);
    setPollData(data);
    setState(State.liveResults);
  }, []);

  useEffect(() => {
    connect(ws, roomId, clientId, 'ws', intialRetryTime)
  }, []);

  useEffect(()=>{
    if(!ws.current) return;

    ws.current.onmessage = (e) =>{
      console.log("event", e)

      const message = JSON.parse(e.data);

      console.log("message", message)

      if(message.error != null)
      {
        setState(State.error)
        setError(message.error)
      }
      else if(message.state != null){
        switch (message.state){
          case 0:
            setState(State.loading)
            break
          case 1:
            console.log("setting question....")
            setQuestion(message.response)
            setState(State.question)
            break
          case 2:
            console.log("setting responses ....")
            setPollData(message.response)
            setState(State.liveResults)
            break
        }
      }
    }

  })

  return [
    state === State.loading && <App key={1} />,
    state === State.question && (
      <AudienceQuestionView
        key={2}
        data={question}
        loading={state === State.loading}
        setState={setStatePollData}
        clientId={clientId}
        roomId={roomId}
      />
    ),
    state === State.liveResults && (
      <AudienceLiveResultsView
        key={3}
        question={question.question}
        options={question.options}
        count={pollData}
        loading={state === State.loading}
      />
    ),
  ];
};

export default AudienceView;
