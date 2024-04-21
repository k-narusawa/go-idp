import React from "react";
import ReactDOM from "react-dom";

const App = () => {
  return (
    <>
      <h1>Hello, React!</h1>
      <p>
        Edit <code>src/index.tsx</code> and save to reload.
      </p>
    </>
  );
};

ReactDOM.render(<App />, document.getElementById("app"));
