import React, { useState, useEffect } from 'react';
import MultipleChoiceQuestionCard from '../Components/MultipleChoiceQuestionCard'
import {Button, Box} from '@mui/material'

const AudienceQuestionView = ({ws}) => {
    const [selected,setSelected] = useState(-1);

    //TODO: need to pass this as a prop
    const clientId = 1;

    let question = "Who is the Captain of Indian Cricket Team";
    let answers = ["Kohli", "Rohit","Pant"];

    const [data, setData] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        if (!ws.current) return;

        ws.onmessage = e => {
            const message = JSON.parse(e.data);
            console.log("e", message);
            setData(message)
            setLoading(false)
        };
    }, []);

    return(
        //TODO: handle loading and error states.
        (!loading) &&
        (<div>
            <MultipleChoiceQuestionCard 
            question={data.question} 
            choices={data.answers}
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
        </div>)
    );

}

export default AudienceQuestionView;