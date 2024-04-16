import { Inter } from "next/font/google";
import { getSession, signIn } from "next-auth/react";

const inter = Inter({ subsets: ["latin"] });

export default function Home() {
  const onLogin = () => {
    console.log("sign in");
    signIn("go-idp");
  }
  return (
    <>
    <div>
      <h1>TOP</h1>
      <button onClick={onLogin}>Login</button>
    </div>
    </>
  );
}
