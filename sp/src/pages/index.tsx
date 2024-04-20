import { getSession, signIn, signOut } from "next-auth/react";
import { GetServerSideProps } from "next";
import { Session } from "next-auth";
import { Button } from "@/components/button";
import axios from "axios";
import { create } from "@github/webauthn-json";
import { parseCreationOptionsFromJSON } from "@github/webauthn-json/browser-ponyfill";

type Props = {
  session: Session | null;
};

const Home = ({ session }: Props) => {
  const onLogin = () => {
    signIn("my-client");
  };

  const onLogout = () => {
    signOut();
  };

  const onPasskey = async () => {
    const json = await axios
      .get("/api/resources/webauthn")
      .then((response) => {
        console.log(response.data);
        return response.data;
      })
      .catch((error) => {
        console.error(error);
        return null;
      });

    if (!json) {
      return;
    }

    try {
      parseCreationOptionsFromJSON(json);
    } catch (e) {
      console.error(e);
      return;
    }

    const response = await create(json);
  };

  if (session) {
    return (
      <>
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
                passkey
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

  return (
    <>
      <div className="p-4">
        <span className="text-2xl font-bold mb-4">TOP</span>
        <div className="flex justify-center">
          <div className="p-4 w-full sm:w-48">
            <Button
              onClick={onLogin}
              variant="primary"
              size="default"
              disabled={false}
            >
              Login
            </Button>
          </div>
        </div>
      </div>
    </>
  );
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
