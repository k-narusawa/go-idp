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
      username: "",
      password: "",
    },
  });

  const onWebauthn = async () => {
    const options = await axios
      .get("/authentication/webauthn/options")
      .then((response) => {
        return response.data;
      })
      .catch((error) => {
        console.error(error);
        return null;
      });

    if (!options) {
      setError("unexpected error occurred.");
      return;
    }

    const parsedOptions = parseRequestOptionsFromJSON({ publicKey: options });

    const credentials = await get(parsedOptions);

    const loginSkipResp = await axios
      .post(`/authentication/webauthn/login`, credentials.toJSON())
      .then((response) => {
        return response.data;
      })
      .catch(() => {
        return null;
      });

    if (!loginSkipResp) {
      setError("unexpected error occurred.");
      return;
    }

    const urlParams = new URLSearchParams(window.location.search);

    const pushTo =
      "/oauth2/session" +
      "?" +
      urlParams.toString() +
      `&token=${loginSkipResp.login_skip_token}`;

    window.location.href = pushTo;

    return;
  };

  return (
    <>
      <div className="p-4">
        <Card>
          <div className="p-4">
            <div className="p-4 flex justify-center text-xl font-semi-bold">
              Login
            </div>
            {error && <div className="p-4 text-red-500">{error}</div>}
            <form method="post">
              <input type="hidden" name="scopes" value="openid" />
              <input type="hidden" name="scopes" value="offline" />
              <div className="p-4">
                <label>Email</label>
                <Input
                  type="text"
                  name="username"
                  placeholder=""
                  control={control}
                />
              </div>
              <div className="p-4">
                <label>Password</label>
                <Input
                  type="password"
                  name="password"
                  placeholder=""
                  control={control}
                />
              </div>
              <div className="p-4 px-12">
                <Button type="submit" variant="primary" disabled={false}>
                  Login
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
                Biometrics Login
              </Button>
            </div>
          </div>
        </Card>
      </div>
    </>
  );
};

ReactDOM.render(<LoginPage />, document.getElementById("app"));
