module.exports = (app) => {
  return async (req, res, next) => {
    const requestedStream = req.params.username,
      endUrl = req.params.endUrl,
      key = `${requestedStream}/${endUrl}`;

    const chunks = [];
    req.on("data", function (chunk) {
      if (endUrl.endsWith(".ts")) app.redisClient.append(key, chunk);
      else chunks.push(chunk);
    });
    req.on("end", function () {
      if (endUrl.endsWith(".ts")) app.redisClient.expire(key, 20);
      else app.redisClient.set(key, Buffer.concat(chunks), { EX: 20 });
      res.status(200).end("ok");
    });
  };
};
