import React from "react";
import { isDoStatement } from "typescript";
import "./App.css";
import ChatHandler from "./chathandler";
import { IMessage } from "./message";

class ChatMessage extends React.Component<IMessage, {}> {
  constructor(props: IMessage) {
    super(props);
  }

  render(): JSX.Element {
    return (
      <div className="message-item">
        <div>
          <b>{this.props.sender}</b>: {this.props.message}
        </div>
      </div>
    );
  }
}

interface ChatProps {
  handler: ChatHandler;
}

interface ChatState {
  input: string;
  messages: JSX.Element[];
}

export class ChatApp extends React.Component<ChatProps, ChatState> {
  bottom: HTMLDivElement | null | undefined;

  constructor(props: ChatProps) {
    super(props);
    this.props.handler.on("message", (message: IMessage) =>
      this.appendMessage(message)
    );

    this.props.handler.on("user-join", (user) =>
      this.appendContent(`${user} has joined!`)
    );
    this.props.handler.on("user-leave", (user) =>
      this.appendContent(`${user} has left!`)
    );
    this.props.handler.on("user-kick", (user) =>
      this.appendContent(`${user} has been kicked!`)
    );

    this.state = {
      input: "",
      messages: [],
    };
  }

  handleSend(content: string) {
    this.props.handler.send(content);
    this.setInput("");
  }

  appendContent(content: string) {
    let state = this.state;
    state.messages.push(
      <div className="content">
        <b>{content}</b>
      </div>
    );
    this.setState(state);
    this.bottom?.scrollIntoView({ behavior: "auto" });
  }

  appendMessage(message: IMessage) {
    let state = this.state;
    state.messages.push(
      <ChatMessage
        message={message.message}
        sender={message.sender}
        key={state.messages.length}
      ></ChatMessage>
    );
    this.setState(state);
    this.bottom?.scrollIntoView({ behavior: "auto" });
  }

  setInput(input: string) {
    this.setState({
      input: input,
      messages: this.state.messages,
    });
  }

  render(): JSX.Element {
    return (
      <>
        <div className="app">
          <div id="message-container">{this.state.messages}</div>
        </div>
        <div id="send-container">
          <form>
            <input
              type="text"
              value={this.state.input}
              onChange={(event) => this.setInput(event.target.value)}
            ></input>
            <button
              type="submit"
              onClick={(event) => {
                event.preventDefault();
                this.handleSend(this.state.input);
              }}
            >
              Send
            </button>
          </form>
        </div>
        <div
          style={{ float: "left", clear: "both" }}
          ref={(el) => {
            this.bottom = el;
          }}
        ></div>
      </>
    );
  }
}

interface ChatRegisterState {
  username: string;
  password: string;
}

export class ChatRegisterPage extends React.Component<
  ChatProps,
  ChatRegisterState
> {
  constructor(props: ChatProps) {
    super(props);

    this.state = {
      username: "",
      password: "",
    };
  }

  handleUsernameChange(event: React.ChangeEvent<HTMLInputElement>) {
    this.setState({
      username: event.target.value,
      password: this.state.password,
    });
  }

  handlePasswordChange(event: React.ChangeEvent<HTMLInputElement>) {
    this.setState({
      username: this.state.username,
      password: event.target.value,
    });
  }

  handleSubmit(event: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    event.preventDefault();
    this.props.handler.login(this.state.username, this.state.password);
    this.setState({ username: "", password: "" });
  }

  render(): JSX.Element {
    return (
      <>
        <form>
          <div>
            Username:{" "}
            <input
              type="text"
              onChange={(e) => this.handleUsernameChange(e)}
              value={this.state.username}
            ></input>
          </div>
          <div>
            Password:{" "}
            <input
              type="password"
              onChange={(e) => this.handlePasswordChange(e)}
              value={this.state.password}
            ></input>
          </div>
          <button type="submit" onClick={(event) => this.handleSubmit(event)}>
            Login
          </button>
        </form>
      </>
    );
  }
}

interface AppState {
  connected: boolean;
  loggedin: boolean;
}

export default class App extends React.Component<{}, AppState> {
  private chathandler: ChatHandler;

  constructor(props: any) {
    super(props);
    this.chathandler = new ChatHandler();
    this.state = { connected: false, loggedin: this.isSocketConnected() };

    this.chathandler.socket.onopen = () => this.setConnected(true);
    this.chathandler.socket.onclose = () => this.setConnected(false);

    this.chathandler.on("connect", () => this.setLoggedIn(true));
    this.chathandler.on("disconnect", () => this.setLoggedIn(false));
  }

  isSocketConnected(): boolean {
    return this.chathandler.socket.readyState === this.chathandler.socket.OPEN;
  }

  setConnected(connected: boolean) {
    this.setState({ connected, loggedin: this.state.loggedin });
  }

  setLoggedIn(loggedin: boolean) {
    this.setState({ loggedin, connected: this.state.connected });
  }

  render(): JSX.Element {
    if (!this.state.loggedin) {
      if (!this.state.connected) {
        return <h1>Socket connecting...</h1>;
      }
      return <ChatRegisterPage handler={this.chathandler} />;
    }
    return <ChatApp handler={this.chathandler} />;
  }
}
