import { Card, Button, CardHeader, Box, Typography } from "@mui/material";

const MultipleChoiceQuestionCard = ({
  question,
  choices,
  selected,
  setSelected,
}) => {
  console.log(question);

  const onClickHandler = (index, selected, setSelected, e) => {
    if (selected === index) {
      setSelected(-1);
    } else {
      setSelected(index);
    }
  };

  return (
    <Card
      sx={{
        transition: 0.3,
        marginBottom: "10px",
        boxShadow: "0 4px 8px 0 rgba(0,0,0,0.2)",
        "&:hover": { boxShadow: "0 8px 16px 0 rgba(0,0,0,0.2)" },
      }}
    >
      <CardHeader title={question} sx={{ alignSelf: "center" }} />
      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          alignItems: "self-start",
        }}
      >
        {choices.map((choice, index) => {
          return (
            <Button
              sx={{ width: 300 }}
              variant={index === selected ? "contained" : "text"}
              onClick={(e) => {
                onClickHandler(index, selected, setSelected, e);
              }}
              key={index}
            >
              <Typography>{choice.option}</Typography>
            </Button>
          );
        })}
      </Box>
    </Card>
  );
};

export default MultipleChoiceQuestionCard;
