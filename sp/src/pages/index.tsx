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

type Props = {
  email: string | null;
  passkeys: PasskeyResponse | null;
  sessionExpired: boolean;
};

const Home = ({ email, passkeys, sessionExpired }: Props) => {
  const [error, setError] = useState<boolean>(false);

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

  if (email) {
    return (
      <>
        <Header onLogout={onLogout} />
        <AccountCard email={email} />

        <div className="p-4" />

        <PasskeyCard
          passkeys={passkeys}
          onRegister={onPasskey}
          onDelete={onDelete}
        />
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
