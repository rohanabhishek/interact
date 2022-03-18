import React, { useState } from 'react';
import MultipleChoiceQuestionCard from './Components/MultipleChoiceQuestionCard'

export default function AudienceView(data){
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
        </div>
    );

}