import { useEffect, useState } from "react";
import LiveResultsComponent from "../Components/LiveResultsComponent";
import socket from "../socketio";

const AudienceLiveResultsView = ()=>{
    let question = "Who is the Captain of Indian Cricket Team";
    let results = [{"option": "kohli","percentage": 20}, {"option": "Rohit","percentage": 50}, {"option": "Pant","percentage": 30}]
    
    //TODO: Loading and error handling
    
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(()=>{
        //for first time call server to add it to LiveResults room
        socket.join("livePollResults")
    },[])

    //start listening to the live poll data
    useEffect(()=>{
        socket.on("livePollData", (data)=>{
            setData(data)
        })
    },[socket])

    return(
        <LiveResultsComponent question={data.question} results={data.results} />
    );
}   

export default AudienceLiveResultsView;