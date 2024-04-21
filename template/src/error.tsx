import React from "react";
import ReactDOM from "react-dom";

const ErrorPage = () => {
  return (
    <>
      <h1>Error</h1>
      <p>
        Edit <code>src/index.tsx</code> and save to reload.
      </p>
    </>
  );
};

ReactDOM.render(<ErrorPage />, document.getElementById("app"));
