import * as React from "react";

import { classNames } from "../../utils";
import { ConnectionHandler } from "../../connection/ws";
import { EventState } from "../../objects/State";

type ChatProps = {
    active: boolean;
    playerID: string;
    conn: ConnectionHandler;
    es: EventState;
}

type ChatState = {
    chats: any[];
}

export class Chat extends React.Component<ChatProps, ChatState> {
    constructor(props: any) {
        super(props)

        this.state = { chats: [] }
        this.props.conn.subscribe('message', this.receiveMessages.bind(this))
    }

    public receiveMessages(chats: any[]) {
        this.setState({ chats })
    }

    public render() {
        return <div className={classNames({"hidden": !this.props.active})}>
            <h3>Game Chat</h3>
            <div>
                {this.state.chats.map((chat) => {
                    return <span>
                        <b>{this.props.es.getPlayer(chat.playerID).username}</b>
                        {chat.message}
                    </span>;
                })}
            </div>
            <div>
                <input type="text" name="chat" placeholder="Type to room" />
            </div>
        </div>;
    }
}