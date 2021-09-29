import EventEmitter from "events";
import { IMessage } from "./message";

const WEBSOCKET_BACKEND: string = `ws://${window.location.host}/ws`

export interface IUserSession {
  username: string;
}

export default class ChatHandler extends EventEmitter {
  socket: WebSocket;
  session: IUserSession | undefined;

  constructor() {
    super();

    this.socket = new WebSocket(WEBSOCKET_BACKEND);

    this.socket.addEventListener("close", () => window.location.reload());

    this.socket.onmessage = (message: any) => {
      const packet = JSON.parse(message["data"]);
      if (!packet || !packet["type"]) return;
      const type: string = packet["type"];

      switch (type) {
        case "login-response":
          if (packet["success"] && packet["username"]) {
            // type = login response
            this.session = {
              username: packet["username"],
            };
            this.emit("connect");
          }
          break;
        case "message":
          if (packet["sender"] && packet["content"]) {
            const message: IMessage = {
              message: packet["content"],
              sender: packet["sender"],
            };
            this.emit("message", message);
          }
          break;
        case "join":
          if (packet["username"]) {
            let username = packet["username"];
            this.emit("user-join", username);
          }
          break;
        case "leave":
          if (packet["username"]) {
            let username = packet["username"];
            this.emit("user-leave", username);
          }
          break;
        case "kick":
          if (packet["username"]) {
            let username = packet["username"];
            this.emit("user-kick", username);
          }
      }
    };
  }

  login(username: string, password: string) {
    this.socket.send(
      JSON.stringify({
        type: "login",
        username: username,
        password: password,
      })
    );
  }

  send(content: string) {
    if (!this.session) return;
    const data = {
      type: "message",
      content: content,
    };
    this.socket.send(JSON.stringify(data));
  }
}
