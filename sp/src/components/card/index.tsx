export const Card = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="mx-auto max-w-md rounded-lg bg-white shadow">
      {children}
    </div>
  );
};
