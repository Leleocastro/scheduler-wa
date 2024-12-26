const functions = require("firebase-functions");
const admin = require("firebase-admin");
const axios = require("axios");

admin.initializeApp();

exports.createUserConsumer = functions.auth
  .user()
  .onCreate(async (user: any) => {
    const username = user.email || user.displayName || user.uid;
    const customId = user.uid;

    const body = {
      username: username,
      custom_id: customId,
    };

    try {
      const response = await axios.post(
        "http://ltag.ddns.net/gateway/consumer",
        body
      );

      if (response.status === 200) {
        console.log("Consumer criado com sucesso no API Gateway");
      } else {
        console.error("Erro ao criar consumer: ", response.statusText);
      }
    } catch (error) {
      console.error("Erro ao chamar API: ", error);
    }
  });
