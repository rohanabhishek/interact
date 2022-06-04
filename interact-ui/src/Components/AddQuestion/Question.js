import { Card } from "@mui/material";
import TextField from "@mui/material/TextField";
import Fab from "@mui/material/Fab";
import AddIcon from "@mui/icons-material/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import React, { useState, useContext } from "react";
import { UserContext } from "../../UserContext.js";
// import { useNavigate } from "react-router-dom";
import "./styles.css";

/*
TODO:
CSS:
Align onto same line the buttons and question text
Differentiablity between question and b/w options
*/

const TextFieldView = (props) => {
  return (
    <TextField
      className="textfieldWidth"
      id="outlined-textarea"
      label={props.label}
      placeholder={props.placeholder}
      value={props.value}
      multiline
      onChange={(e) => props.onChange(e)}
      margin="normal"
    />
  );
};

const OptionTextView = (props) => {
  return (
    <div>
      <div className="text-child">
        <TextFieldView
          label={props.label}
          placeholder={props.placeholder}
          value={props.value}
          onChange={(e) => props.onChange(e)}
        />
        <Fab className="fab" size="small" color="secondary" aria-label="delete">
          <DeleteIcon onClick={() => props.handleDeleteOption()} />
        </Fab>
      </div>
    </div>
  );
};

const QuestionCard = ({ setState }) => {
  const { contextDetails } = useContext(UserContext);
  let roomId = contextDetails.roomId;
  // let navigate = useNavigate();
  // TODO: Retrieve and use UserId(Host Id).

  const [question, setQuestion] = useState("");
  const [options, updateOptions] = useState([]);

  const handleAddOption = () => {
    let newOptions = [...options, ""];
    updateOptions(newOptions);
  };

  const handleDeleteOption = (index) => {
    let newOptions = [...options];
    newOptions.splice(index, 1);
    updateOptions(newOptions);
  };

  const handleUpdateOption = (value, index) => {

    //option object
    let option = {
      idx: index,
      option: value
    }

    let newOptions = [...options];
    newOptions[index] = option;
    updateOptions(newOptions);
  };

  const handlePostQuestion = () => {
    let postQuestionURI =
      process.env.REACT_APP_SERVER_ADDR + roomId + `/addLiveQuestion`;
    fetch(postQuestionURI, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      // TODO: Add a form to collect and send this body
      // Add a dropdown select like TextField for different questionType.
      // Can use id instead of owner explicit value.
      body: JSON.stringify({
        owner: "IAmTheHost",
        questionType: "Single MCQ",
        question: question,
        options: options,
      }),
    })
      .then((response) => {
        console.log(response);
        if (!response.ok) {
          console.log("response not ok");
        } else {
          let questionData = {
            question: question,
            options: options,
          };
          setState(questionData);
          // navigate("/AudienceView", { replace: true });
        }
        return response.json();
      })
      .then((data) => {
        console.log("response json:", data);
      });
  };

  return (
    <div className="div-style">
      <Card
        className="center-display"
        variant="outlined"
        sx={{
          transition: 0.3,
          marginBottom: "10px",
          boxShadow: "0 4px 8px 0 rgba(0,0,0,0.2)",
          "&:hover": { boxShadow: "0 8px 16px 0 rgba(0,0,0,0.2)" },
        }}
      >
        <div>
          <div className="text-child">
            <TextFieldView
              value={question}
              label={"Question"}
              placeholder={"Add Question content here"}
              onChange={(e) => setQuestion(e.target.value)}
            />
            <Fab className="fab" size="small" color="primary" aria-label="add">
              <AddIcon onClick={() => handleAddOption()} />
            </Fab>
          </div>
        </div>
        <div>
          {options.map((option, index) => {
            return (
              <OptionTextView
                label={"Option " + (index + 1).toString()}
                placeholder={
                  "Add Option " + (index + 1).toString() + "content here"
                }
                value={option.option}
                key={index}
                onChange={(e) => handleUpdateOption(e.target.value, index)}
                handleDeleteOption={(i) => handleDeleteOption(i)}
              />
            );
          })}
        </div>
        <Fab
          variant="extended"
          className="fab"
          size="medium"
          aria-label="edit"
          color="primary"
          onClick={() => handlePostQuestion()}
        >
          SEND
        </Fab>
      </Card>
    </div>
  );
};

export default QuestionCard;
