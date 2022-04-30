import { useEffect, useState, useRef } from "react";
import LiveResultsComponent from "../Components/LiveResultsComponent";

const AudienceLiveResultsView = ()=>{
    // let question = "Who is the Captain of Indian Cricket Team";
    // let results = [{"option": "kohli","percentage": 20}, {"option": "Rohit","percentage": 50}, {"option": "Pant","percentage": 30}]

    //TODO: Loading and error handling
    const ws = useRef(null);

    const[data, setData] = useState(null)
    const[loading, setLoading] = useState(true)

    useEffect(() => {
        ws.current = new WebSocket("ws://localhost:8080/1/liveResults");
        ws.current.onopen = () => console.log("ws opened");
        ws.current.onclose = () => console.log("ws closed");

        const wsCurrent = ws.current;

        return () => {
            wsCurrent.close();
        };
    }, []);

    useEffect(() => {
        if (!ws.current) return;

        ws.current.onmessage = e => {
            const message = JSON.parse(e.data);
            console.log("e", message);
            setData(message)
            setLoading(false)
        };
    }, []);

    //TODO: Handle loading and error cases
    return(
        (!loading && <LiveResultsComponent question={data.question} results={data.results} />)
    );
}

export default AudienceLiveResultsView;