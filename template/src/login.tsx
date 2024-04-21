import React from "react";
import ReactDOM from "react-dom";

const LoginPage = () => {
  return (
    <>
      <h1>Login</h1>
      <form method="post">
        <input type="hidden" name="scopes" value="openid" />
        <input type="hidden" name="scopes" value="offline" />
        <div className="form-control">
          <label>Username</label>
          <input
            type="text"
            id="username"
            name="username"
            value="test@example.com"
          />
        </div>
        <div className="form-control">
          <label>Password</label>
          <input
            type="password"
            id="password"
            name="password"
            value="!Password0"
          />
        </div>
        <button type="submit">ログイン</button>
      </form>
    </>
  );
};

ReactDOM.render(<LoginPage />, document.getElementById("app"));
