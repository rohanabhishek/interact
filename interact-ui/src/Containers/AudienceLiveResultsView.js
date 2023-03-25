import { useEffect, useState, useRef } from "react";
import LiveResultsComponent from "../Components/LiveResultsComponent";

const AudienceLiveResultsView = ({ options, loading, question, count }) => {

  //TODO: Handle loading and error cases

  return (
    !loading && (
      <LiveResultsComponent
        question={question}
        options={options}
        count={count}
      />
    )
  );
};

export default AudienceLiveResultsView;
