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

  if (req.method === "GET") {
    const apiResponse = await axios
      .get(
        `${process.env.IDP_URL}/resources/users/registrations/webauthn/options`,
        {
          headers: {
            Authorization: `Bearer ${session.accessToken}`,
          },
        }
      )
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
    return;
  } else if (req.method === "POST") {
    const apiResponse = await axios
      .post(
        `${process.env.IDP_URL}/resources/users/registrations/webauthn/result`,
        req.body,
        {
          headers: {
            Authorization: `Bearer ${session.accessToken}`,
          },
          params: {
            challenge: req.query.challenge,
          },
        }
      )
      .then((response) => response.data)
      .catch((error) => {
        console.error(error);
        return null;
      });

    res.status(200).json(apiResponse);
    return;
  } else {
    res.status(405).end("Method Not Allowed");
    return;
  }
}
