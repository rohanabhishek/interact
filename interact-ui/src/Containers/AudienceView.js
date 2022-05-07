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

const State = {
  loading: 0,
  question: 1,
  results: 2,
  error: 3,
};

const AudienceView = () => {
  const { contextDetails } = useContext(UserContext);
  let roomId = contextDetails.roomId;
  let clientId = contextDetails.clientId;

  const wsResults = useRef(null);
  const wsQuestion = useRef(null);

  const [state, setState] = useState(State.loading);
  const [pollData, setPollData] = useState(null);
  const [question, setQuestion] = useState(null);

  const setStatePollData = useCallback((data) => {
    console.log(data);
    setPollData(data);
    setState(State.results);
  }, []);

  useEffect(() => {
    connect(wsQuestion, roomId, clientId, `liveQuestion`);
    connect(wsResults, roomId, clientId, `liveResults`);
  }, []);

  useEffect(() => {
    if (!wsQuestion.current) return;

    wsQuestion.current.onmessage = (e) => {
      console.log(e);
      const message = JSON.parse(e.data);
      console.log("e", message);
      if (message.stateChange != null) {
        setState(message.stateChange);
      } else {
        console.log(message.options);
        setQuestion(message);
        setState(State.question);
      }
    };
  }, []);

  useEffect(() => {
    if (!wsResults.current) return;

    wsResults.current.onmessage = (e) => {
      console.log(e);
      const message = JSON.parse(e.data);
      console.log("e", message);
      if (message.stateChange != null) {
        setState(message.stateChange);
      } else {
        setPollData(message);
      }
    };
  }, []);

  return [
    state === state.loading && <App key={1} />,
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
    state === State.results && (
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

function connect(ws, roomId, clientId, path) {
  ws.current = new WebSocket(
    `ws://localhost:8080/${roomId}/${path}/${clientId}`
  );

  console.log(ws);

  //TODO: hanlde it correctly
  ws.current.onclose = function (e) {
    console.log(
      "Socket is closed. Reconnect will be attempted in 1 second.",
      e.reason
    );
    setTimeout(function () {
      connect(ws, roomId, clientId, path);
    }, 1000);
  };

  ws.current.onerror = function (err) {
    console.error("Socket encountered error: ", err.message, "Closing socket");
    ws.current.close();
  };
}

export default AudienceView;
