import * as React from "react";

import { classNames } from "../../utils";
import { ConnectionHandler } from "../../connection/ws";
import { EventState } from "../../objects/State";
import { Message } from "./chat/Message";

type ChatProps = {
    active: boolean;
    playerID: string;
    conn: ConnectionHandler;
    es: EventState;
}

type ChatState = {
    chats: any[];
    text?: string;
}

export class Chat extends React.Component<ChatProps, ChatState> {
    constructor(props: any) {
        super(props)

        this.state = { chats: [], text: null }
        this.props.conn.subscribe('message', this.receiveMessages.bind(this))
    }

    public receiveMessages(msg: any) {
        this.setState({ chats: msg.data })
    }

    public textUpdate(event: any) {
        this.setState({text: event.target.value});
    }

    public submit(event: any) {
        event.preventDefault();

        let action = {
            type: 'message',
            data: {
                message: this.state.text,
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
                {this.state.chats.map((chat) => {
                    return <Message player={this.props.es.getPlayer(chat.player_id)} message={chat.message}/>;
                })}
            </div>
            <div>
                <form onSubmit={this.submit.bind(this)}>
                    <input type="text" name="chat" placeholder="Type to room"
                        value={this.state.text}
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