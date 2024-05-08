import { getSession, signIn, signOut } from "next-auth/react";
import { GetServerSideProps } from "next";
import axios from "axios";
import {
  create,
  parseCreationOptionsFromJSON,
} from "@github/webauthn-json/browser-ponyfill";
import { useEffect, useState } from "react";
import { AccountCard } from "@/components/pages/top/AccountCard";
import { PasskeyCard } from "@/components/pages/top/PasskeyCard";
import { Header } from "@/components/common/Header";
import { SessionCard } from "@/components/pages/top/SessionCard";
import { Toast } from "@/components/common/Toast";

type Props = {
  accessToken: string;
  refreshToken: string;
  email: string | null;
  passkeys: PasskeyResponse | null;
  sessionExpired: boolean;
};

const Home = ({
  accessToken,
  refreshToken,
  email,
  passkeys,
  sessionExpired,
}: Props) => {
  const [error, setError] = useState<boolean>(false);
  const [copyComplete, setCopyComplete] = useState<boolean>(false);

  useEffect(() => {
    if (sessionExpired) {
      onLogout();
    }

    if (!email && !passkeys) {
      signIn("go-idp");
    }
  }, [email, passkeys, sessionExpired]);

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
      .catch((error) => {
        console.error(error);
        setError(true);
      });
    if (error) {
      return;
    }
    window.location.reload();
  };

  const onDelete = async (id: string) => {
    await axios
      .delete(`/api/resources/users/webauthn`, {
        params: {
          id: id,
        },
      })
      .then((response) => response.data)
      .catch((error) => {
        console.error(error);
        return null;
      });

    window.location.reload();
  };

  const copyAccessToken = async () => {
    await navigator.clipboard.writeText(accessToken);
    setCopyComplete(true);
  };

  const copyRefreshToken = async () => {
    await navigator.clipboard.writeText(refreshToken);
    setCopyComplete(true);
  };

  if (email) {
    return (
      <>
        <Header onLogout={onLogout} />
        <AccountCard email={email} />

        <div className="p-4" />

        <SessionCard
          copyAccessToken={copyAccessToken}
          copyRefreshToken={copyRefreshToken}
        />

        <div className="p-4" />

        <PasskeyCard
          passkeys={passkeys}
          onRegister={onPasskey}
          onDelete={onDelete}
        />

        {copyComplete && (
          <div className="absolute inset-x-0 bottom-0 flex items-center justify-center">
            <Toast
              message="Copied to clipboard"
              type="success"
              onClose={() => setCopyComplete(false)}
            />
          </div>
        )}
      </>
    );
  }

  return <></>;
};

export const getServerSideProps: GetServerSideProps = async (context) => {
  const session = await getSession({ req: context.req });
  console.log(session?.accessToken);
  console.log(session?.refreshToken);

  if (!session) {
    return {
      props: {
        sessionExpired: true,
      },
    };
  }

  try {
    const resp = await axios.get(
      `${process.env.IDP_URL}/resources/users/webauthn`,
      {
        headers: {
          Authorization: `Bearer ${session.accessToken}`,
        },
      }
    );
    return {
      props: {
        accessToken: session.accessToken,
        refreshToken: session.refreshToken,
        email: session.email,
        passkeys: resp.data,
        sessionExpired: false,
      },
    };
  } catch (error) {
    return {
      props: {
        sessionExpired: true,
      },
    };
  }
};

export default Home;
