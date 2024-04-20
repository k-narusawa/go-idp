import { authOptions } from "@/pages/api/auth/[...nextauth]";
import axios from "axios";
import { NextApiRequest, NextApiResponse } from "next";
import { getServerSession } from "next-auth";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const session = await getServerSession(req, res, authOptions);

  if (!session) {
    res.status(401).end("Unauthorized");
    return;
  }

  const apiResponse = await axios
    .get(`${process.env.IDP_URL}/api/v1/users/webauthn`, {
      headers: {
        Authorization: `Bearer ${session.accessToken}`,
      },
    })
    .then((response) => response.data)
    .catch((error) => {
      console.error(error);
      return null;
    });

  if (!apiResponse) {
    res.status(500).end("Internal Server Error");
    return;
  }

  res.status(200).json(apiResponse);
}
