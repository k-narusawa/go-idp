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

  if (req.method === "DELETE") {
    const id = req.query.id as string;
    await axios
      .delete(`${process.env.IDP_URL}/resources/users/webauthn/${id}`, {
        headers: {
          Authorization: `Bearer ${session.accessToken}`,
        },
        params: {
          challenge: req.query.challenge,
        },
      })
      .then((response) => response.data)
      .catch((error) => {
        console.error(error);
        return null;
      });

    res.status(204).json(undefined);
    return;
  }

  res.status(405).end("Method Not Allowed");
  return;
}
