type Props = {
  onLogout: () => void;
};

export const Header: React.FC<Props> = ({ onLogout }) => {
  return (
    <>
      <nav className="bg-white border-gray-200 dark:bg-gray-900">
        <div className="flex flex-wrap justify-between items-center mx-auto max-w-screen-xl p-4">
          <span className="self-center text-2xl font-semibold whitespace-nowrap dark:text-white">
            Go-IdP
          </span>
          <div className="flex items-center space-x-6 rtl:space-x-reverse">
            <div
              className="text-sm text-blue-600 dark:text-blue-500 hover:underline hover:cursor-pointer"
              onClick={onLogout}
            >
              Logout
            </div>
          </div>
        </div>
      </nav>
    </>
  );
};
