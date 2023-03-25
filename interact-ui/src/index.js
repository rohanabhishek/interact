import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import AudienceView from "./Containers/AudienceView";
import reportWebVitals from "./reportWebVitals";
import StartEventPage from "./Containers/StartEventPage/StartEventPage";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { UserContextProvider } from "./UserContext.js";
import HostView from "./Containers/HostView";
import Homepage from "./Containers/Homepage";
// const roomId = "3aedf06b-9170-4e99-adc1-10d6126b756a"
// const clientId = "06cefd8f-fba7-4eca-9059-4dc90fb071d5"

ReactDOM.render(
  <UserContextProvider>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Homepage />} />
        <Route path="hostView" element={<HostView />} />
        <Route path="AudienceView" element={<AudienceView />} />
      </Routes>
    </BrowserRouter>
  </UserContextProvider>,
  document.getElementById("root")
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
