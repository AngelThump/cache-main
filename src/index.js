const config = require("../config/config.json");
const { createClient } = require("redis");
const express = require("express");
const app = express();
app.disable("x-powered-by");
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
    //password: config.redis.password,
  });

  await app.redisClient
    .connect()
    .then(() => {
      console.info("Redis client connected.");
    })
    .catch((e) => console.error(e));
});
const cache = require("./cache");
const auth = require("./auth");

app.post("/hls/:username/:endUrl", auth(app), cache(app));
