const express = require("express");
const bodyParser = require("body-parser");
const app = express();
const WebSocket = require("ws");
const wss = new WebSocket.Server({ noServer: true });
const fs = require("fs");

const { Client, LocalAuth } = require("whatsapp-web.js");

app.use(bodyParser.json());
let sessions = {};

const whatsappRouter = express.Router();

app.use("/whatsapp", whatsappRouter);

whatsappRouter.post("/create-session", (req, res) => {
  const { phoneNumber } = req.body;

  if (sessions[phoneNumber]?.client) {
    console.log(`Session already exists for: ${phoneNumber}`);
    return res.status(400).send({
      message: "Session already exists",
    });
  }

  sessions[phoneNumber] = { status: "pending" };

  console.log("Creating session for: ", phoneNumber);

  const authStrategy = new LocalAuth({
    clientId: phoneNumber,
  });

  const client = new Client({
    authStrategy: authStrategy,
    takeoverOnConflict: true,
    takeoverTimeoutMs: 10,
    puppeteer: {
      headless: true,
      args: [
        "--no-sandbox",
        "--disable-dev-shm-usage",
        "--disabled-setupid-sandbox",
        "--disable-accelerated-2d-canvas",
        "--no-first-run",
        "--no-zygote",
        "--single-process", // Adicione este argumento se necessário
        "--disable-gpu",
      ],
    },
  });

  sessions[phoneNumber].client = client;

  client.on("qr", (qr) => {
    console.log("QR CODE RECEIVED:", qr);

    sessions[phoneNumber].qrCode = qr;
  });

  client.once("ready", () => {
    console.log("Client is ready!");

    sessions[phoneNumber].isReady = true;
  });

  client.on("authenticated", () => {
      console.log("Authenticated");
  });
  
  client.on("auth_failure", (msg) => {
      console.error("AUTH FAILURE", msg);
  });
  
  client.on("disconnected", (reason) => {
      console.log("Client was logged out", reason);
  });

  client
    .initialize()
    .then(() => {
      console.log(`WHATSAPP WEB SESSION STARTED FOR ${phoneNumber}`);
      sessions[phoneNumber].status = "connected";
    })
    .catch((error) => {
      console.error(`Failed to initialize session for ${phoneNumber}:`, error);
      delete sessions[phoneNumber];
    });

  res.send({
    message: "Session created",
  });
});

whatsappRouter.post("/restart-session", async (req, res) => {
  const { phoneNumber } = req.body;

  if (!sessions[phoneNumber]?.client) {
    console.log(`No session found for: ${phoneNumber}`);
    return res.status(400).send({
      message: "No session found",
    });
  }

  client.on("qr", (qr) => {
    console.log("QR CODE RECEIVED:", qr);

    sessions[phoneNumber].qrCode = qr;
  });

  client.once("ready", () => {
    console.log("Client is ready!");

    sessions[phoneNumber].isReady = true;
  });

  console.log(`Restarting session for: ${phoneNumber}`);

  await sessions[phoneNumber].client.destroy();

  sessions[phoneNumber].client.initialize().then(() => {
    console.log(`WHATSAPP WEB SESSION RESTARTED FOR ${phoneNumber}`);
    sessions[phoneNumber].status = "connected";
  });

  res.send({
    message: "Session restarted",
  });
});

whatsappRouter.post("/close-session", (req, res) => {
  const { phoneNumber } = req.body;

  if (!sessions[phoneNumber]?.client) {
    console.log(`No session found for: ${phoneNumber}`);
    return res.status(400).send({
      message: "No session found",
    });
  }

  console.log(`Closing session for: ${phoneNumber}`);

  sessions[phoneNumber].client.destroy().then(() => {
    console.log(`WHATSAPP WEB SESSION CLOSED FOR ${phoneNumber}`);
    delete sessions[phoneNumber];
  });

  res.send({
    message: "Session closed",
  });
});

whatsappRouter.get("/get-qrcode/:phoneNumber", (req, res) => {
  const { phoneNumber } = req.params;

  if (!sessions[phoneNumber]?.client) {
    console.log(`No session found for: ${phoneNumber}`);
    return res.status(400).send({
      message: "No session found",
    });
  }

  if (sessions[phoneNumber]?.isReady) {
    console.log(`Client is already ready for: ${phoneNumber}`);
    return res.status(400).send({
      message: "Client is Ready, you can send messages",
    });
  }

  if (!sessions[phoneNumber]?.qrCode) {
    console.log(`No QR code found for: ${phoneNumber}`);
    return res.status(400).send({
      message: "No QR code found, try again",
    });
  }

  res.send({
    qrCode: sessions[phoneNumber].qrCode,
  });
});

whatsappRouter.post("/send-message", async (req, res) => {
  const { phoneNumber, message, to } = req.body;

  if (!sessions[phoneNumber]?.client) {
    console.log(`No session found for: ${phoneNumber}`);
    return res.status(400).send({
      message: "No session found",
    });
  }

  const chatId = to + "@c.us";

  try {
    let response = await sessions[phoneNumber].client.sendMessage(
      chatId,
      message
    );
    console.log("Message sent: ", response);
    res.send({
      message: "Message sent",
    });
  } catch (err) {
    console.error("Failed to send message:", err);
    res.status(500).send({
      message: "Failed to send message",
    });
  }
});

wss.on("connection", (ws, req) => {
  const phoneNumber = req.url.split("/")[2];

  if (!sessions[phoneNumber]?.client) {
    console.log("No session found for: ", phoneNumber);
    ws.send(
      JSON.stringify({
        type: "error",
        message: "No session found",
      })
    );
    ws.close();
    return;
  }

  let client = sessions[phoneNumber].client;

  setInterval(() => {
    if (ws.readyState === WebSocket.OPEN) {
      ws.ping();
    }
  }, 30000);

  client.on("qr", (qr) => {
    console.log("QR CODE RECEIVED:", qr);
    ws.send(
      JSON.stringify({
        type: "qr",
        qrCode: qr,
      })
    );
  });

  client.once("ready", () => {
    console.log("Client is ready!");
    ws.send(
      JSON.stringify({
        type: "ready",
        message: "Client is ready",
      })
    );
  });

  client.on("error", (err) => {
    console.error("Client error:", err);
    ws.send(
      JSON.stringify({
        type: "error",
        message: err.message,
      })
    );
  });

  client.on("message", (message) => {
    console.log("MESSAGE RECEIVED:", message);
    ws.send(
      JSON.stringify({
        type: "message",
        message: message,
      })
    );
  });

  client.on("loading_screen", (status) => {
    console.log("LOADING SCREEN:", status);
    ws.send(
      JSON.stringify({
        type: "loading_screen",
        status: status,
      })
    );
  });

  ws.on("message", async (message) => {
    const data = JSON.parse(message);

    if (data.type === "message") {
      try {
        const chatId = data.phoneNumber + "@c.us";

        let response = await client.sendMessage(chatId, data.message);
        console.log("Message sent: ", response);
      } catch (err) {
        ws.send(
          JSON.stringify({
            type: "error",
            message: err.message,
          })
        );
      }
    }
  });
});

const server = app.listen(3000, () => {
  console.log("Server started on port 3000");
});

server.on("upgrade", (request, socket, head) => {
  const pathname = request.url.split("/")[1];

  if (pathname === "ws") {
    wss.handleUpgrade(request, socket, head, (ws) => {
      wss.emit("connection", ws, request);
    });
  } else {
    socket.destroy();
  }
});
