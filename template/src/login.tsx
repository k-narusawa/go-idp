import React from "react";
import axios from "axios";
import ReactDOM from "react-dom";
import { SubmitHandler, useForm } from "react-hook-form";
import "./styles.css";
import { Button } from "./components/Button";
import { Card } from "./components/Card";
import { Input } from "./components/Input";
import { HorizontalLine } from "./components/HorizontalLine";
import {
  get,
  parseRequestOptionsFromJSON,
} from "@github/webauthn-json/browser-ponyfill";

const LoginPage = () => {
  type Inputs = {
    username: string;
    password: string;
  };

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<Inputs>({
    defaultValues: {
      username: "test@example.com",
      password: "!Password0",
    },
  });

  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    var queryParams = window.location.search;
    queryParams = queryParams + "&";

    const params = new URLSearchParams();
    params.append("username", data.username);
    params.append("password", data.password);
    params.append("scopes", "openid");
    params.append("scopes", "offline");

    const body = {
      username: data.username,
      password: data.password,
    };

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

  const onWebauthn = async () => {
    const options = await axios
      .get("/api/v1/webauthn/login")
      .then((response) => {
        return response.data;
      })
      .catch((error) => {
        console.error(error);
        return null;
      });

    if (!options) {
      console.error("WebAuthn login failed");
      return;
    }

    const parsedOptions = parseRequestOptionsFromJSON({
      publicKey: options.publicKey,
    });

    console.log(parsedOptions);

    const credentials = await get(parsedOptions);
  };

  return (
    <>
      <div className="p-4">
        <Card>
          <div className="p-4">
            <div className="p-4 flex justify-center text-xl font-semi-bold">
              ログイン
            </div>
            <form onSubmit={handleSubmit(onSubmit)}>
              <div className="p-4">
                <label>ログインID</label>
                <Input
                  type="text"
                  name="username"
                  placeholder="test@example.com"
                  control={control}
                  rules={{
                    required: "メールアドレスの入力は必須です",
                  }}
                />
              </div>
              <div className="p-4">
                <label>パスワード</label>
                <Input
                  type="password"
                  name="password"
                  placeholder="********"
                  control={control}
                  rules={{
                    required: "パスワードの入力は必須です",
                  }}
                />
              </div>
              <div className="p-4 px-12">
                <Button type="submit" variant="primary" disabled={false}>
                  ログイン
                </Button>
              </div>
            </form>
            <HorizontalLine />
            <div className="pt-4 px-12">
              <Button
                type="button"
                variant="primary"
                disabled={false}
                onClick={onWebauthn}
              >
                生体認証ログイン
              </Button>
            </div>
          </div>
        </Card>
      </div>
    </>
  );
};

function bufferDecode(value: string) {
  return Uint8Array.from(
    atob(value.replace(/-/g, "+").replace(/_/g, "/")),
    (c) => c.charCodeAt(0)
  );
}

ReactDOM.render(<LoginPage />, document.getElementById("app"));
