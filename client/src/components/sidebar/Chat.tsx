import * as React from "react";

import { classNames } from "../../utils";
import { ConnectionHandler, EventType } from "../../connection/ws";
import { Message } from "./chat/Message";
import { AnonymousPlayer, findPlayer } from "../../objects/Player";

type ChatMessage = {
    player_id: number;
    player_username: string;
    message: string;
};

type ChatProps = {
    active: boolean;
    playerID: number;
    players: AnonymousPlayer[];
    conn: ConnectionHandler;
}

type ChatState = {
    chats: ChatMessage[];
    text?: string;
}

export class Chat extends React.Component<ChatProps, ChatState> {
    constructor(props: any) {
        super(props)

        this.state = { chats: [], text: null }
        this.props.conn.subscribe(EventType.message, this.receiveMessages.bind(this))
    }

    public receiveMessages(msg: any) {
        this.setState({ chats: msg.data })
    }

    public textUpdate(event: any) {
        this.setState({text: event.target.value});
    }

    public submit(event: any) {
        event.preventDefault();

        const player = findPlayer(this.props.playerID, this.props.players);
        const action = {
            type: 'message',
            data: {
                message: this.state.text,
                username: player?.username,
            }
        };

        this.props.conn.send(JSON.stringify(action))
        this.setState({text: ''})

        return false;
    }

    public render() {
        return <div className={classNames({"hidden": !this.props.active})}>
            <h3>Game Chat</h3>
            <div>
                {this.state.chats.map((chat: ChatMessage, i: number) =>
                    <Message key={i} playerUsername={chat.player_username} message={chat.message}/>
                )}
            </div>
            <div>
                <form onSubmit={this.submit.bind(this)}>
                    <input type="text" name="chat" placeholder="Type to room"
                        value={this.state.text || ""}
                        onChange={this.textUpdate.bind(this)}/>
                    <button disabled={!this.state.text} type="submit"
                        className={classNames("btn-large grey darken-3", {'disabled': !this.state.text})}>
                        Submit
                    </button>
                </form>
            </div>
        </div>;
    }
}