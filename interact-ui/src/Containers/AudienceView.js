import React, { useState, useEffect, useRef } from 'react';
import AudienceLiveResultsView from './AudienceLiveResultsView';
import AudienceQuestionView from './AudienceQuestionView';
import App from '../App';

const State = {
    loading : 0,
    question : 1,
    results : 2
}

const  AudienceView = ({roomId,clientId})=>{
    
    const liveQuestionSocket = useRef(null);
    const liveResultsSocket = useRef(null);

    const [loading, setLoading] = useState(false)
    const [state, setState] = useState(State.loading)

    useEffect(()=>{
        connect(`${roomId}/liveResults`, liveResultsSocket)
        connect(`${roomId}/liveQuestion`, liveQuestionSocket)
    },[])

    return(
        <div>
            <App/>
            {state == State.question} && <AudienceQuestionView ws={{liveQuestionSocket}}/>
            {state == State.results} && <AudienceLiveResultsView ws = {{liveResultsSocket}}/>
        </div>
    )
}

function connect(socket, ws) {
    ws.current = new WebSocket(`ws://localhost:8080/${socket}`);

    console.log(ws)
 
   ws.current.onclose = function(e) {
     console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
     setTimeout(function() {
       connect(socket,ws);
     }, 1000);
   };
 
   ws.current.onerror = function(err) {
     console.error('Socket encountered error: ', err.message, 'Closing socket');
     ws.current.close();
   };
 }

export default AudienceView