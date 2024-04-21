import React from "react";
import axios from "axios";
import ReactDOM from "react-dom";
import { SubmitHandler, useForm } from "react-hook-form";
import { register } from "ts-node";

const LoginPage = () => {
  type Inputs = {
    username: string;
    password: string;
  };

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Inputs>();

  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    const queryParams = window.location.search;

    const params = new URLSearchParams();
    params.append("username", data.username);
    params.append("password", data.password);
    params.append("scopes", "openid");
    params.append("scopes", "offline");

    const res = await axios
      .post(`/oauth2/auth${queryParams}`, params, {
        withCredentials: true,
      })
      .then((response) => {
        return response.data;
      })
      .catch((error) => {
        console.error(error);
        return null;
      });

    if (res) {
      window.location.href = res.redirect_to;
    } else {
      console.error("Login failed");
    }
  };

  return (
    <>
      <div className="login-container">
        <h1>ログイン</h1>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="form-control">
            <label> Username</label>
            <input type="text" {...register("username", { required: true })} />
            {errors.username && <span>Username is required</span>}
          </div>
          <div className="form-control">
            <label>Password</label>
            <input
              type="password"
              {...register("password", { required: true })}
            />
            {errors.password && <span>Password is required</span>}
          </div>
          <button type="submit">ログイン</button>
        </form>
      </div>
    </>
  );
};

ReactDOM.render(<LoginPage />, document.getElementById("app"));
