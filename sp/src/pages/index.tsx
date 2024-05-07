import { getSession, signIn, signOut } from "next-auth/react";
import { GetServerSideProps } from "next";
import { Session } from "next-auth";
import { Button } from "@/components/common/Button";
import axios from "axios";
import {
  create,
  parseCreationOptionsFromJSON,
} from "@github/webauthn-json/browser-ponyfill";
import { Toast } from "@/components/common/Toast";
import { useEffect, useState } from "react";
import { Card } from "@/components/common/Card";
import { HorizontalLine } from "@/components/common/HorizontalLine";
import { profile } from "console";
import { AccountCard } from "@/components/pages/top/AccountCard";

type Props = {
  session: Session | null;
};

const Home = ({ session }: Props) => {
  const [success, setSuccess] = useState<boolean>(false);
  const [error, setError] = useState<boolean>(false);

  useEffect(() => {
    if (!session) {
      signIn("go-idp");
    }
  }, [session]);

  const onLogout = () => {
    signOut({ callbackUrl: "signOut", redirect: true });
  };

  const onPasskey = async () => {
    const options = await axios
      .get("/api/resources/users/registrations/webauthn/options")
      .then((response) => {
        return response.data;
      })
      .catch((error) => {
        console.error(error);
        return null;
      });

    const challenge = options.challenge;
    if (!options) {
      return;
    }

    const parsedOptions = parseCreationOptionsFromJSON({ publicKey: options });
    const response = await create(parsedOptions);
    await axios
      .post(
        "/api/resources/users/registrations/webauthn/result",
        response.toJSON(),
        {
          params: {
            challenge: challenge,
          },
        }
      )
      .then(() => {
        setSuccess(true);
      })
      .catch((error) => {
        console.error(error);
        setError(true);
      });
  };

  if (session) {
    return (
      <>
        {success && (
          <Toast
            message="パスキー登録成功"
            type="success"
            onClose={() => setSuccess(false)}
          />
        )}
        {error && (
          <Toast
            message="パスキー登録失敗"
            type="danger"
            onClose={() => setError(false)}
          />
        )}
        <div className="p-4">
          <span className="text-2xl font-bold mb-4">TOP</span>

          <AccountCard email={session.email} />

          <div className="flex justify-center">
            <div className="p-4 w-full sm:w-48">
              <Button
                onClick={onPasskey}
                variant="primary"
                size="default"
                disabled={false}
              >
                パスキー登録
              </Button>
            </div>
          </div>
          <div className="flex justify-center">
            <div className="p-4 w-full sm:w-48">
              <Button
                onClick={onLogout}
                variant="danger"
                size="default"
                disabled={false}
              >
                Logout
              </Button>
            </div>
          </div>
        </div>
      </>
    );
  }

  return <></>;
};

export const getServerSideProps: GetServerSideProps = async (context) => {
  const session = await getSession({ req: context.req });
  console.log(session?.accessToken);
  console.log(session?.refreshToken);

  return {
    props: {
      session,
    },
  };
};

export default Home;
