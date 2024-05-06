import { Card } from "@/components/common/Card";
import { HorizontalLine } from "@/components/common/HorizontalLine";

type Props = {
  email: string;
};

export const AccountCard: React.FC<Props> = ({ email }) => {
  return (
    <>
      <Card>
        <div className="p-4 flex justify-center text-2xl font-semi-bold">
          Account
        </div>
        <HorizontalLine />
        <div className="grid grid-cols-6 py-4 px-8">
          <div className="col-start-1">
            <span className="text-gray-500">Email</span>
          </div>
          <div className="col-start-5 whitespace-nowrap">{email}</div>
        </div>
        <HorizontalLine />
        <div className="grid grid-cols-6 py-4 px-8">
          <div className="col-start-1">
            <span className="text-gray-500">Password</span>
          </div>
          <div className="col-start-5 whitespace-nowrap">********</div>
        </div>
      </Card>
    </>
  );
};
