import { getSession, signIn, signOut } from "next-auth/react";
import { GetServerSideProps } from "next";
import { Session } from "next-auth";
import { Button } from "@/components/Button";
import axios from "axios";
import {
  create,
  parseCreationOptionsFromJSON,
} from "@github/webauthn-json/browser-ponyfill";
import { Toast } from "@/components/Toast";
import { useEffect, useState } from "react";

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
      .get("/api/resources/webauthn")
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

    console.log(options);

    const parsedOptions = parseCreationOptionsFromJSON({ publicKey: options });

    const response = await create(parsedOptions);

    console.log(response.toJSON());

    await axios
      .post("/api/resources/webauthn", response.toJSON(), {
        params: {
          challenge: challenge,
        },
      })
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
        <div className="p-4 overflow-auto">
          <span className="text-2xl font-bold mb-4">TOP</span>
          <table className="w-full table-auto">
            <tbody>
              <tr className="border-b border-gray-200">
                <td className="py-2 px-4">id</td>
                <td className="py-2 px-4 whitespace-normal">{session.id}</td>
              </tr>
              <tr className="border-b border-gray-200">
                <td className="py-2 px-4">accessToken</td>
                <td className="py-2 px-4 whitespace-normal">
                  {session.accessToken}
                </td>
              </tr>
              <tr className="border-b border-gray-200">
                <td className="py-2 px-4">refreshToken</td>
                <td className="py-2 px-4 whitespace-normal">
                  {session.refreshToken}
                </td>
              </tr>
              <tr className="border-b border-gray-200">
                <td className="py-2 px-4">idToken</td>
                <td className="py-2 px-4 whitespace-normal max-w-[200px] overflow-x-auto">
                  {session.idToken}
                </td>
              </tr>
              <tr className="border-b border-gray-200">
                <td className="py-2 px-4">expiresAt</td>
                <td className="py-2 px-4 whitespace-normal">
                  {session.expiresAt}
                </td>
              </tr>
            </tbody>
          </table>
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

  return {
    props: {
      session,
    },
  };
};

export default Home;
