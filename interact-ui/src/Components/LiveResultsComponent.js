import { Typography , Card, CardHeader, Box, LinearProgress, CardContent} from "@mui/material";
import PropTypes from 'prop-types';

const LiveResultsComponent = ({question, results}) => {
    return(

        <Box alignContent={'center'} flex={1}>
          <Card
            variant="outlined"
            sx={{
                height: "100vh",
                width: "50%",
                margin:'auto',
                transition: 0.3,
                marginTop: "10px",
                marginBottom: "10px",
                boxShadow: "0 4px 8px 0 rgba(0,0,0,0.2)",
                "&:hover": {boxShadow: "0 8px 16px 0 rgba(0,0,0,0.2)"}
            }}
          >
          <CardHeader title={question} sx={{alignSelf: 'center'}} />

          <Box sx={{ display: 'flex', flexDirection: 'column' }}>

              {results.map((result,i) =>{
                return(
                  <Box>
                    <CardContent key={i}>
                      <Typography gutterBottom variant="h5" component="div">
                        {result.option}
                      </Typography>
                      <LinearProgressWithLabel value={result.percentage}/>
                    </CardContent>
                  </Box>
              )})}
          </Box>
          </Card>
        </Box>

    );
}


function LinearProgressWithLabel(props) {
    return (
      <Box sx={{ display: 'flex', alignItems: 'center' }}>
        <Box sx={{ width: '75%', mr: 1 }}>
          <LinearProgress sx={{height: "15px", borderTopRightRadius: "5px", borderBottomRightRadius: "5px"}} variant="determinate" value={props.value} />
        </Box>
        <Box >
          <Typography variant="h6">{`${Math.round(
            props.value,
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