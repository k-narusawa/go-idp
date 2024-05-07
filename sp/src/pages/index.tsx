import { getSession, signIn, signOut } from "next-auth/react";
import { GetServerSideProps } from "next";
import { Button } from "@/components/common/Button";
import axios from "axios";
import {
  create,
  parseCreationOptionsFromJSON,
} from "@github/webauthn-json/browser-ponyfill";
import { Toast } from "@/components/common/Toast";
import { useEffect, useState } from "react";
import { AccountCard } from "@/components/pages/top/AccountCard";
import { PasskeyCard } from "@/components/pages/top/PasskeyCard";

type Props = {
  email: string | null;
  passkeys: PasskeyResponse | null;
};

const Home = ({ email, passkeys }: Props) => {
  const [error, setError] = useState<boolean>(false);

  useEffect(() => {
    if (!email) {
      signIn("go-idp");
    }
  }, [email]);

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
        <div className="p-4">
          <span className="text-2xl font-bold mb-4">TOP</span>

          <AccountCard email={email} />

          <div className="p-4" />

          <PasskeyCard
            passkeys={passkeys}
            onRegister={onPasskey}
            onDelete={onDelete}
          />

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

  const apiResponse = await axios
    .get(`${process.env.IDP_URL}/resources/users/webauthn`, {
      headers: {
        Authorization: `Bearer ${session?.accessToken}`,
      },
    })
    .then((response) => response.data)
    .catch((error) => {
      console.error(error);
      return null;
    });

  if (!apiResponse) {
    return {
      props: {
        session: null,
      },
    };
  }

  return {
    props: {
      email: session?.email,
      passkeys: apiResponse,
    },
  };
};

export default Home;
