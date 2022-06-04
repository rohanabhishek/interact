import { 
    Card, 
    CardContent, 
    CardHeader,  
    Stack, 
    Tooltip,
    Container
} from "@mui/material"
import LinearProgress, { linearProgressClasses } from '@mui/material/LinearProgress';
import { styled } from '@mui/material/styles';


const LiveResultsView = ({question, options, results})=>{

    console.log("results", results)
    console.log("options", options)
    return (  
       <Container>
            <Card>
                <CardHeader
                    title={question}
                    titleTypographyProps={{variant:'h4' }}
                />
                <CardContent>
                    <Stack>
                        {options.map((x)=>(
                                <ResultsComponent 
                                    result={results === null? null : results.find(y=> y.idx == x.idx)} 
                                    option={x}
                                />
                            )
                        )}
                    </Stack>
                </CardContent>
            </Card>
        </Container>
    )
}

const BorderLinearProgress = styled(LinearProgress)(({ theme }) => ({
    height: 10,
    borderRadius: 5,
    [`&.${linearProgressClasses.colorPrimary}`]: {
      backgroundColor: 'transparent',
    },
    [`& .${linearProgressClasses.bar}`]: {
      borderRadius: 5,
      backgroundColor: theme.palette.mode === 'light' ? '#1a90ff' : '#308fe8',
    },
  }));


const ResultsComponent = ({result, option})=>{

    console.log("result", result)
    console.log("option", option)

    return(
        <Card 
            key={option.idx}
            sx={{
                backgroundColor: 'transparent',
                boxShadow: 'none'
        }}>
            <CardHeader 
                title={option.option}
                titleTypographyProps={{variant:'h5'}}
            />
            <CardContent>
                <Tooltip title={`${result!= null ? result.count: 0} responses`} placement="bottom-start">
                    <BorderLinearProgress variant="determinate" value={parseInt(result!= null ? result.percentage: 0)} />
                </Tooltip> 
            </CardContent>
        </Card>       
    )
}

export default LiveResultsView