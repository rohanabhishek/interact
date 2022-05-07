import { Card } from "@mui/material";
import Fab from "@mui/material/Fab";
import Alert from "@mui/material/Alert";
import React, { useState, useRef, useContext } from "react";
import PlayArrowOutlinedIcon from "@mui/icons-material/PlayArrowOutlined";
import PersonAddAltOutlinedIcon from "@mui/icons-material/PersonAddAltOutlined";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import { UserContext } from "../../UserContext.js";
import { useNavigate } from "react-router-dom";
import "./styles.css";

// TODO: integrate to existing UI code.
const StartEventPage = () => {
  let navigate = useNavigate();
  const [openDialogBox, setOpen] = useState(false);
  const eventIDRef = useRef("");
  const [createEventError, setCreateEventError] = useState("");
  const [joinEventError, setJoinEventError] = useState("");
  const { setContextDetails } = useContext(UserContext);

  const createEventHandler = () => {
    let createEventURI = process.env.REACT_APP_SERVER_ADDR + "createEvent";
    fetch(createEventURI, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      // TODO: Add a form to collect and send this body
      body: JSON.stringify({
        hostName: "IAmTheHost",
        eventType: "LivePolls",
      }),
    })
      .then((response) => {
        if (!response.ok) {
          console.log("response not ok");
          setCreateEventError("Join Event failed: API response not ok");
        }
        return response.json();
      })
      .then((data) => {
        if (data.roomId) {
          // TODO: Add hostId also to context
          setContextDetails({ roomId: data.roomId });
          navigate("/AddQuestion", { replace: true });
        } else {
          setCreateEventError("Create event falied: API didn't return roomId");
        }
        console.log("response json:", data);
      })
      .catch((error) => {
        console.log("error: ", error);
        setCreateEventError(error.toString());
      });
  };
  const joinEventHandler = (eventID) => {
    let joinEventURI =
      process.env.REACT_APP_SERVER_ADDR + eventID + "/joinEvent";
    fetch(joinEventURI, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((response) => {
        console.log(response);
        if (!response.ok) {
          console.log("response not ok");
        }
        return response.json();
      })
      .then((data) => {
        if (data.clientId) {
          setContextDetails({
            roomId: eventIDRef.current.value,
            clientId: data.clientId,
          });
          navigate("/AudienceView", { replace: true });
        }
        console.log("response json:", data);
        setOpen(false);
      })
      .catch((error) => {
        if (!error.response) {
          setJoinEventError(error.toString());
        } else {
          setJoinEventError(error.response.data.message);
        }
      });
  };

  const closeDialogHandler = () => {
    setOpen(false);
  };
  return (
    <div className="div-style">
      {createEventError && (
        <Alert variant="outlined" severity="error">
          {createEventError}
        </Alert>
      )}
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
        <div className="div-style-block">
          <Fab
            variant="extended"
            size="medium"
            className="fab-block"
            onClick={(e) => createEventHandler()}
          >
            <PlayArrowOutlinedIcon sx={{ mr: 5 }} />
            <b>START EVENT</b>
          </Fab>
          <Fab
            variant="extended"
            size="medium"
            className="fab-block"
            onClick={(e) => {
              setOpen(true);
            }}
          >
            <PersonAddAltOutlinedIcon sx={{ mr: 5 }} />
            <b> JOIN EVENT</b>
          </Fab>
        </div>
      </Card>
      <Dialog open={openDialogBox} onClose={closeDialogHandler}>
        <DialogTitle>Join Event</DialogTitle>
        {joinEventError && (
          <Alert variant="outlined" severity="error">
            {joinEventError}
          </Alert>
        )}
        <DialogContent>
          <DialogContentText>
            To join event, type in the event-id shared by the HOST.
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="name"
            label="event-id"
            type="text"
            fullWidth
            variant="standard"
            inputRef={eventIDRef}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={closeDialogHandler}>Cancel</Button>
          <Button onClick={(e) => joinEventHandler(eventIDRef.current.value)}>
            POST
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default StartEventPage;
