import React, { useEffect } from "react";
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
  const [error, setError] = React.useState<string | null>(null);
  useEffect(() => {
    if (window.hasOwnProperty("idpMessage") && (window as any).idpMessage) {
      setError((window as any).idpMessage);
    }
  }, []);

  type Inputs = {
    username: string;
    password: string;
  };

  const {
    control,
    formState: { errors },
  } = useForm<Inputs>({
    defaultValues: {
      username: "test@example.com",
      password: "!Password0",
    },
  });

  const onWebauthn = async () => {
    const options = await axios
      .get("/webauthn/login")
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

    const credentials = await get(parsedOptions);

    await axios
      .post(`/webauthn/login`, credentials.toJSON())
      .then((response) => {
        console.log(response);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  return (
    <>
      <div className="p-4">
        <Card>
          <div className="p-4">
            <div className="p-4 flex justify-center text-xl font-semi-bold">
              ログイン
            </div>
            {error && <div className="p-4 text-red-500">{error}</div>}
            <form method="post">
              <input type="hidden" name="scopes" value="openid" />
              <input type="hidden" name="scopes" value="offline" />
              <div className="p-4">
                <label>ログインID</label>
                <Input
                  type="text"
                  name="username"
                  placeholder="test@example.com"
                  control={control}
                />
              </div>
              <div className="p-4">
                <label>パスワード</label>
                <Input
                  type="password"
                  name="password"
                  placeholder="********"
                  control={control}
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
