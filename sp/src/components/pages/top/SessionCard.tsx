import { Card } from "@/components/common/Card";
import { HorizontalLine } from "@/components/common/HorizontalLine";
import { ClipBoardIcon } from "@/components/common/Icon";

type Props = {
  copyAccessToken: () => Promise<void>;
  copyRefreshToken: () => Promise<void>;
};

export const SessionCard: React.FC<Props> = ({
  copyAccessToken,
  copyRefreshToken,
}) => {
  return (
    <>
      <Card>
        <div className="p-4 flex justify-center text-2xl font-semi-bold">
          Session
        </div>
        <HorizontalLine />
        <div className="grid grid-cols-6 py-4 px-8">
          <div className="col-start-1">
            <span className="text-gray-500">AccessToken</span>
          </div>
          <div className="col-start-7 whitespace-nowrap">
            <div className="hover:cursor-pointer" onClick={copyAccessToken}>
              <ClipBoardIcon />
            </div>
          </div>
        </div>
        <HorizontalLine />
        <div className="grid grid-cols-6 py-4 px-8">
          <div className="col-start-1">
            <span className="text-gray-500">RefreshToken</span>
          </div>
          <div className="col-start-7 whitespace-nowrap">
            <div className="hover:cursor-pointer" onClick={copyRefreshToken}>
              <ClipBoardIcon />
            </div>
          </div>
        </div>
      </Card>
    </>
  );
};
