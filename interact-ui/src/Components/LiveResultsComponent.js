import {
  Typography,
  Card,
  CardHeader,
  Box,
  LinearProgress,
} from "@mui/material";
import PropTypes from "prop-types";

const LiveResultsComponent = ({ question, options, count }) => {
  return (
    <Card
      variant="outlined"
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
        {options.map((option, i) => {
          return (
            <Box sx={{ flexDirection: " column", width: "25%" }} key={i}>
              <Typography>{option}</Typography>
              <LinearProgressWithLabel
                value={count ? (option in count ? count[option] : 0) : 0}
              />
            </Box>
          );
        })}
      </Box>
    </Card>
  );
};

function LinearProgressWithLabel(props) {
  return (
    <Box sx={{ display: "flex", alignItems: "center" }}>
      <Box sx={{ width: "100%", mr: 1 }}>
        <LinearProgress
          sx={{ height: "10px", borderRadius: 5 }}
          variant="determinate"
          {...props}
        />
      </Box>
      <Box sx={{ minWidth: 35 }}>
        <Typography variant="body2" color="text.secondary">{`${Math.round(
          props.value
        )}%`}</Typography>
      </Box>
    </Box>
  );
}

LinearProgressWithLabel.propTypes = {
  /**
   * The value of the progress indicator for the determinate and buffer variants.
   * Value between 0 and 100.
   */
  value: PropTypes.number.isRequired,
};

export default LiveResultsComponent;
