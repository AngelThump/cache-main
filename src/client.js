const config = require("../config/config.json");
const feathers = require("@feathersjs/feathers");
const socketio = require("@feathersjs/socketio-client");
const io = require("socket.io-client");

module.exports = () => {
  const socket = io(config.socket.hostname, {
    extraHeaders: {
      "X-Api-Key": config.socket.authKey,
    },
  });
  
  const client = feathers();

  client.configure(
    socketio(socket, {
      timeout: 2000,
    })
  );
  return client;
};
