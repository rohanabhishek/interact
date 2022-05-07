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
  question: 0,
  results: 1,
  error: 2,
  end: 3,
};

const HostView = () => {
  const { contextDetails } = useContext(UserContext);
  let roomId = contextDetails.roomId;
  let clientId = contextDetails.userId;

  const wsResults = useRef(null);

  const [state, setState] = useState(State.question);
  const [pollData, setPollData] = useState(null);
  const [question, setQuestion] = useState(null);

  const setStateQuestion = useCallback((data) => {
    console.log(data);
    setQuestion(data);
    setState(State.results);
  }, []);

  const changeStateToQuestion = useCallback(() => {
    setState(State.question);
  }, []);

  useEffect(() => {
    connect(wsResults, roomId, clientId, `liveResults`);
  }, []);

  useEffect(() => {
    if (!wsResults.current) return;

    wsResults.current.onmessage = (e) => {
      console.log(e);
      const message = JSON.parse(e.data);
      console.log("e", message);
      if (message.stateChange != null) {
        // TODO: Handle host separately server-side, not to send this to host
      } else {
        setPollData(message);
      }
    };
  }, []);

  return [
    state === State.question && <QuestionCard setState={setStateQuestion} />,
    state === State.results && (
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

export default HostView;
