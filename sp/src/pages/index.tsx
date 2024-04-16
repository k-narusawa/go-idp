import { getSession, signIn, signOut } from "next-auth/react";
import { GetServerSideProps } from "next";
import { Session } from "next-auth";

type Props = {
  session: Session | null;
};

const Home = ({ session }: Props) => {
  const onLogin = () => {
    signIn("my-client");
  }

  const onLogout = () => {
    signOut();
  }

  if(session) {
    return (
      <>
        <div className="p-4 overflow-auto">
          <h1 className="text-2xl font-bold mb-4">TOP</h1>
          <table className="w-full table-auto">
            <tbody>
              <tr className="border-b border-gray-200">
                <td className="py-2 px-4">id</td>
                <td className="py-2 px-4 whitespace-normal">{session.id}</td>
              </tr>
              <tr className="border-b border-gray-200">
                <td className="py-2 px-4">accessToken</td>
                <td className="py-2 px-4 whitespace-normal">{session.accessToken}</td>
              </tr>
              <tr className="border-b border-gray-200">
                <td className="py-2 px-4">refreshToken</td>
                <td className="py-2 px-4 whitespace-normal">{session.refreshToken}</td>
              </tr>
              <tr className="border-b border-gray-200">
                <td className="py-2 px-4">idToken</td>
                <td className="py-2 px-4 whitespace-normal max-w-[200px] overflow-x-auto">{session.idToken}</td>
              </tr>
              <tr className="border-b border-gray-200">
                <td className="py-2 px-4">expiresAt</td>
                <td className="py-2 px-4 whitespace-normal">{session.expiresAt}</td>
              </tr>
            </tbody>
          </table>
          <div className="p-4"/>
          <button type="button" 
            className="py-2.5 px-5 me-2 mb-2 
            text-sm font-medium text-gray-900 focus:outline-none 
            bg-white rounded-lg border border-gray-200 
            hover:bg-gray-100 hover:text-blue-700 
            focus:z-10 focus:ring-4 focus:ring-gray-100 
            dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700"
            onClick={onLogout}
            >
            Logout
          </button>
        </div>
      </>
    );
  }

  return (
    <>
    <div className="p-4">
      <h1>TOP</h1>
      <div className="p-4"/>
      <button type="button" 
        className="py-2.5 px-5 me-2 mb-2 
        text-sm font-medium text-gray-900 focus:outline-none 
        bg-white rounded-lg border border-gray-200 
        hover:bg-gray-100 hover:text-blue-700 
        focus:z-10 focus:ring-4 focus:ring-gray-100 
        dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700"
        onClick={onLogin}
        >
        Login
      </button>
    </div>
    </>
  );
}

export const getServerSideProps: GetServerSideProps = async (context) => {
  const session = await getSession({ req: context.req });
  return {
    props: {
      session,
    },
  };
}

export default Home;