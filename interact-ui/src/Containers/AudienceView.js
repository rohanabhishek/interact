import React, { useState } from 'react';
import MultipleChoiceQuestionCard from '../Components/MultipleChoiceQuestionCard'
import {Button, Box} from '@mui/material'

const AudienceView = () => {
    const [selected,setSelected] = useState(-1);

    let question = "Who is the Captain of Indian Cricket Team";
    let answers = ["Kohli", "Rohit","Pant"];

    return(
        <div>
        <MultipleChoiceQuestionCard 
            question={question} 
            choices={answers}
            selected = {selected}
            setSelected = {setSelected}
        />

        <MultipleChoiceQuestionCard 
            question={question} 
            choices={answers}
            selected = {selected}
            setSelected = {setSelected}
        />

        <MultipleChoiceQuestionCard 
            question={question} 
            choices={answers}
            selected = {selected}
            setSelected = {setSelected}
        />

        <MultipleChoiceQuestionCard 
            question={question} 
            choices={answers}
            selected = {selected}
            setSelected = {setSelected}
        />
        <Box textAlign='center'>
            <Button
                color='primary'
                size='large'
                type='submit'
                variant='contained'
            >
                Submit
            </Button>
        </Box>

        </div>
    );

}

export default AudienceView;