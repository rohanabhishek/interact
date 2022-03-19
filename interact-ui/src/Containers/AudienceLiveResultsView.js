import LiveResultsComponent from "../Components/LiveResultsComponent";

const AudienceLiveResultsView = ()=>{
    let question = "Who is the Captain of Indian Cricket Team";
    let results = [{"option": "kohli","percentage": 20}, {"option": "Rohit","percentage": 50}, {"option": "Pant","percentage": 30}]

    return(
        <LiveResultsComponent question={question} results={results} />
    );
}   

export default AudienceLiveResultsView;