import React, { useState} from 'react';
import MultipleChoiceQuestionCard from '../Components/MultipleChoiceQuestionCard'
import {Button, Box} from '@mui/material'

const AudienceQuestionView = ({data, loading, setState, clientId, roomId}) => {
    const [selected,setSelected] = useState(-1);

    const onSubmitHandler = ()=>{
        //send response and setState
        if (selected === -1) return

        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                clientId:  clientId,
                response:  data.options[selected]
            })
        };
        fetch(`http://localhost:8080/${roomId}/sendResponse/${clientId}`, requestOptions)
            .then(response => response.json())
            .then(data => {
                console.log(data)
                console.log(data.liveResults)
                setState(data.liveResults)
            });
    }
    
    return(
        //TODO: handle loading and error states.
        (!loading) &&
        (<div>
            <MultipleChoiceQuestionCard 
            question= {data.question} 
            choices={data.options}
            selected = {selected}
            setSelected = {setSelected}
            />
            <Box textAlign='center'>
                <Button
                    color='primary'
                    size='large'
                    type='submit'
                    variant='contained'
                    disabled = {selected === -1}
                    onClick = {onSubmitHandler}
                >
                    Submit
                </Button>
            </Box>
        </div>)
    );
}


export default AudienceQuestionView;