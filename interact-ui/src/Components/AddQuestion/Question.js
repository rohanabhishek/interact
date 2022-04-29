import { Card, CardHeader, Dialog } from "@mui/material";
import Box from "@mui/material/Box";
import Fab from "@mui/material/Fab";
import AddIcon from "@mui/icons-material/Add";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
import TextField from "@mui/material/TextField";
import React from "react";
import "./styles.css";

/*
TODO:
Add API to send the question
CSS:
Align onto same line the buttons and question text
Differentiablity between question and b/w options
*/

const OptionView = (props) => {
  return (
    <div>
      <div className="text-child">
        <Box
          component="div"
          sx={{
            border: "1px solid grey",
            borderRadius: 5,
            borderColor: "black",
            fontWeight: "bold",
            fontSize: 18,
          }}
        >
          <CardHeader
            className="box-left-align"
            title={props.text}
            sx={{ display: "flex", fontSize: 30 }}
          />
        </Box>
      </div>
      <div className="button-child">
        <Fab className="fab" size="small" color="secondary" aria-label="delete">
          <DeleteIcon onClick={() => props.onClickDeleteOption()} />
        </Fab>
        <Fab className="fab" size="small" aria-label="edit">
          <EditIcon onClick={() => props.onClickEditOption()} />
        </Fab>
      </div>
    </div>
  );
};

const OptionsView = (props) => {
  const renderOption = (i) => {
    console.log("renderOption: ", i, props.options[i]);
    return (
      <OptionView
        text={props.options[i]}
        onClickEditOption={() => {
          props.onClickEditOption(i);
        }}
        onClickDeleteOption={() => props.onClickDeleteOption(i)}
      />
    );
  };

  return (
    <div>
      {props.options.map((text, index) => {
        return renderOption(index);
      })}
    </div>
  );
};

class Question extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      questionType: null,
      question: "Add your Question...",
      options: [],
      dialogDefaultText: null,
      isDialogOpen: false,
      editingQuestion: false,
      currentEditOption: { editing: false },
      numOfOptions: 0,
    };
  }

  handleEditOptionClick = (index) => {
    this.setState({ currentEditOption: { editing: true, index: index } });
    let optionText = this.state.options[index];
    this.setState({ dialogDefaultText: optionText });
    this.setState({ isDialogOpen: true });
  };

  handleDeleteOptionClick = (index) => {
    let copyOptions = this.state.options;
    copyOptions.splice(index, 1);
    this.setState({ options: copyOptions });
  };

  addOption = () => {
    let copyOptions = this.state.options;
    copyOptions.concat("");
    this.setState({
      numOfOptions: this.state.numOfOptions + 1,
      options: copyOptions,
    });
    this.handleEditOptionClick(this.state.numOfOptions);
  };

  modifyQuestion = () => {
    this.setState({
      editingQuestion: true,
      isDialogOpen: true,
      dialogDefaultText: this.state.question,
    });
  };

  closeDialogBox = () => {
    this.setState({ isDialogOpen: false });
    if (this.state.currentEditOption.editing) {
      this.setState({ currentEditOption: { editing: false } });
      console.log(this.state.options);
    }
    if (this.state.editingQuestion) {
      this.setState({ editingQuestion: false });
    }
  };

  updateDialogTextChange = (event) => {
    if (this.state.editingQuestion) {
      this.setState({ question: event.target.value });
    }
    if (this.state.currentEditOption.editing) {
      let copyOptions = this.state.options;
      copyOptions[this.state.currentEditOption.index] = event.target.value;
      this.setState({ options: copyOptions });
    }
  };

  render() {
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
              <Box
                component="div"
                sx={{
                  border: "1px solid grey",
                  borderRadius: 5,
                  borderColor: "black",
                  fontWeight: "bold",
                  fontSize: 18,
                }}
              >
                <CardHeader
                  className="box-left-align"
                  title={this.state.question}
                  sx={{ display: "flex", fontSize: 30 }}
                />
              </Box>
            </div>
            <div className="button-child">
              <Fab
                className="fab"
                size="small"
                color="primary"
                aria-label="add"
              >
                <AddIcon onClick={this.addOption} />
              </Fab>
              <Fab className="fab" size="small" aria-label="edit">
                <EditIcon onClick={this.modifyQuestion} />
              </Fab>
            </div>
          </div>
          <div>
            <OptionsView
              options={this.state.options}
              onClickEditOption={(i) => this.handleEditOptionClick(i)}
              onClickDeleteOption={(i) => this.handleDeleteOptionClick(i)}
            />
          </div>
          <Fab
            variant="extended"
            className="fab"
            size="medium"
            aria-label="edit"
            color="primary"
          >
            SEND
          </Fab>
        </Card>

        <Dialog open={this.state.isDialogOpen} onClose={this.closeDialogBox}>
          <TextField
            autoFocus
            margin="dense"
            id="dialogText"
            fullWidth
            variant="standard"
            defaultValue={this.state.dialogDefaultText}
            onChange={this.updateDialogTextChange}
          />
        </Dialog>
      </div>
    );
  }
}

export default Question;
