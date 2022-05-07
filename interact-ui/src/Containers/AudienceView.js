import React, { useState, useEffect } from "react";
import MultipleChoiceQuestionCard from "../Components/MultipleChoiceQuestionCard";
import { Button, Box } from "@mui/material";

const AudienceView = () => {
  const [selected, setSelected] = useState(-1);

  //TODO: need to pass this as a prop
  const clientId = 1;

//   let question = "Who is the Captain of Indian Cricket Team";
//   let answers = ["Kohli", "Rohit", "Pant"];

  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch(`http://127.0.0.1:8080/room/fetchLiveQuestion`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(`The status is ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        console.log("response json:", data);
        setData(data);
        setLoading(false);
      })
      .catch((error) => {
        setError(error);
      });
  }, []);

  if (loading) {
    return (
      <div>
        <h1> Loading... </h1>
      </div>
    );
  } else {
    return (
      //TODO: handle loading and error states.
      <div>
        <MultipleChoiceQuestionCard
          question={data.question}
          choices={data.options}
          selected={selected}
          setSelected={setSelected}
        />
        <Box textAlign="center">
          <Button
            color="primary"
            size="large"
            type="submit"
            variant="contained"
          >
            Submit
          </Button>
        </Box>
      </div>
    );
  }
};

export default AudienceView;
