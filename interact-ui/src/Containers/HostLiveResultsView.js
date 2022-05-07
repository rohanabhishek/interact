import AudienceLiveResultsView from "./AudienceLiveResultsView";
import Fab from "@mui/material/Fab";

const HostLiveResultsView = ({
  question,
  options,
  count,
  roomId,
  changeStateToQuestion,
}) => {
  return (
    <div>
      <AudienceLiveResultsView
        key={3}
        question={question}
        options={options}
        count={count}
        loading={0}
      />
      <Fab
        variant="extended"
        className="fab"
        size="medium"
        aria-label="edit"
        color="primary"
        onClick={() =>
          handleMoveToNextQuestion({ roomId, changeStateToQuestion })
        }
      >
        NEXT
      </Fab>
    </div>
  );
};

const handleMoveToNextQuestion = ({ roomId, changeStateToQuestion }) => {
  let moveToNextQuestionURI =
    process.env.REACT_APP_SERVER_ADDR + roomId + "/nextLiveQuestion";
  fetch(moveToNextQuestionURI, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((response) => {
      if (!response.ok) {
        console.log("response not ok");
      }
      return response.json();
    })
    .then((data) => {
      console.log("response json:", data);
      // Navigate back to addquestion
      changeStateToQuestion();
    })
    .catch((error) => {
      console.error(error);
    });
};

export default HostLiveResultsView;
