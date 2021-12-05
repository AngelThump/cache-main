module.exports = (app) => {
  return async (req, res, next) => {
    const username = req.params.username,
      endUrl = req.params.endUrl;

    const base64String = await app.redisClient.get(username).catch(() => null);
    if (!base64String) return res.status(500).json({ error: true, msg: "Could not find base64 string..." });

    const key = `${base64String}_${username}/${endUrl}`;

    const chunks = [];
    req.on("data", function (chunk) {
      /* Maybe use this when LL-HLS is needed
      if (endUrl.endsWith(".ts")) app.redisClient.append(key, chunk);
      else chunks.push(chunk);*/
      chunks.push(chunk);
    });
    req.on("end", function () {
      /*if (endUrl.endsWith(".ts")) app.redisClient.expire(key, 20);
        else app.redisClient.set(key, Buffer.concat(chunks), { EX: 20 });
      */
      app.redisClient.set(key, Buffer.concat(chunks), { EX: 20 });
      res.status(200).end("ok");
    });
  };
};
