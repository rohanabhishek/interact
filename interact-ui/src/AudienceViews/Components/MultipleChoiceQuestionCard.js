import Card from '@mui/material/Card';
import Button from '@mui/material/Button';
import CardHeader from '@mui/material/CardHeader';
import Box from '@mui/material/Box';





export default function MultipleChoiceQuestionCard({question, choices, selected, setSelected}) {

    const onClickHandler = (index,selected,setSelected,e)=>{
        if( selected === index){
            setSelected(-1);
        }
        else{
            setSelected(index);
        }
    }


    return (
        <Card>
            <CardHeader title={question} />
                <Box sx={{ display: 'flex', flexDirection: 'column', alignItems:'self-start' }}>
                    {choices.map((choice,index) =>{
                       return( 
                            <Button
                               variant= {(index == selected)? "contained" : "text"}
                               onClick={(e)=>{onClickHandler(index,selected,setSelected,e)}}                               
                               >                              
                               {choice}
                            </Button>
                    )})}                             
                </Box>            
        </Card>
    )
}
