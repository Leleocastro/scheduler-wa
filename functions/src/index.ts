import { UserRecord } from "firebase-admin/auth";

import * as functions from "firebase-functions/v1";
import * as admin from "firebase-admin";
import axios from "axios";

admin.initializeApp();

export const createUserConsumer = functions.auth
  .user()
  .onCreate(async (user: UserRecord) => {
    functions.logger.info("User: ", user);
    const username = user.email;
    const customId = user.uid;

    const body = {
      username: username,
      custom_id: customId,
    };

    functions.logger.info("Body: ", body);

    try {
      const response = await axios.post(
        "http://ltag.ddns.net/gateway/consumer",
        body,
        {
          headers: {
            apikey: "VFj2U7wsGCWax04zzwZrSa2vNhAtDMKy",
          },
        }
      );

      if (response.status === 200) {
        functions.logger.info("Consumer criado com sucesso no API Gateway");
      } else {
        functions.logger.error("Erro ao criar consumer: ", response.statusText);
      }
    } catch (error) {
      functions.logger.error("Erro ao chamar API: ", error);
    }
  });
