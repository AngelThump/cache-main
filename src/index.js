const config = require("../config/config.json");
const { createClient } = require("redis");
const express = require("express");
const client = require("./client");
const app = express();
app.disable("x-powered-by");

const sleep = (ms) => {
  return new Promise((resolve) => setTimeout(resolve, ms));
};

app.listen(config.port, async () => {
  console.log(`Angelthump REDIS MAIN listening on port ${config.port}!`);

  app.redisClient = createClient({
    socket: config.redis.useUnixSocket
      ? {
          path: config.redis.unix,
        }
      : {
          host: config.redis.hostname,
        },
    password: config.redis.password,
  });

  await app.redisClient
    .connect()
    .then(() => {
      console.info("Redis client connected.");
      app.redisClient.flushAll();
    })
    .catch((e) => console.error(e));

  app.client = client();

  while (!app.client.io.connected) await sleep(500);

  const streamService = app.client.service("streams");

  const streams = await streamService.find().catch((e) => {
    console.error(e);
    return null;
  });

  for (let stream of streams) {
    app.redisClient.set(stream.username, Buffer.from(stream.createdAt + stream.username).toString("base64"));
  }

  streamService.on("created", (stream) => {
    app.redisClient.set(stream.username, Buffer.from(stream.createdAt + stream.username).toString("base64"));
  });

  streamService.on("removed", (stream) => {
    app.redisClient.del(stream.username);
  });
});
const cache = require("./cache");
const auth = require("./auth");

app.post("/hls/:username/:endUrl", auth(app), cache(app));
