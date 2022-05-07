import { useState, createContext, useMemo } from "react";

const UserContext = createContext();

const UserContextProvider = (props) => {
  const [contextDetails, setContextDetails] = useState({});

  const value = useMemo(
    () => ({ contextDetails, setContextDetails }),
    [contextDetails]
  );

  return (
    <UserContext.Provider value={value}>{props.children}</UserContext.Provider>
  );
};

export { UserContext, UserContextProvider };
