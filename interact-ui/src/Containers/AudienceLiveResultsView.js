import { useEffect, useState, useRef } from "react";
import LiveResultsComponent from "../Components/LiveResultsComponent";

const AudienceLiveResultsView = ({options, loading, question, count})=>{
    // let question = "Who is the Captain of Indian Cricket Team";
    // let results = [{"option": "kohli","percentage": 20}, {"option": "Rohit","percentage": 50}, {"option": "Pant","percentage": 30}]
    
    //TODO: Loading and error handling
    //const ws = useRef(null);

    // const[data, setData] = useState(null)
    // const[loading, setLoading] = useState(true)

    // useEffect(() => {
    //     if (!ws.current) return;

    //     ws.current.onmessage = e => {
    //         const message = JSON.parse(e.data);
    //         console.log("e", message);
    //         setData(message)
    //         setLoading(false)
    //     };
    // }, []);

    //TODO: Handle loading and error cases

    console.log(count)
    return(
        (!loading && <LiveResultsComponent question={question} options={options} count={count} />)
    );
}   

export default AudienceLiveResultsView;