const config = require("../config/config.json");

module.exports = (app) => {
  return async (req, res, next) => {
    if (!req.headers["authorization"]) return res.status(400).json({ error: true, msg: "Missing Auth Key!" });

    const key = req.headers.authorization.split(" ")[1];

    if (key !== config.authKey) return res.status(403).json({ error: true, msg: "Auth Key does not match!" });

    next();
  };
};
