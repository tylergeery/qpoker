import * as React from "react";

import { classNames } from "../../utils";
import { ConnectionHandler } from "../../connection/ws";

type ChatProps = {
    active: boolean;
    playerID: string;
    conn: ConnectionHandler;
}

type ChatState = {
    chats: any[];
}

export class Chat extends React.Component<ChatProps, ChatState> {
    constructor(props: any) {
        super(props)
        this.state = { chats: [] }

        // TODO: register for new chats
    }

    public render() {
        return <div className={classNames({"hidden": !this.props.active})}>
            <h3>Game Chat</h3>
            <div>
                <label>
                    Start Game:
                    <button type="button">
                        Start
                    </button>
                </label>
            </div>
            <h3>Chip Requests</h3>
            <div>
                <label>
                    request:
                    <button type="button">Approve</button>
                    <button type="button">Deny</button>
                </label>
            </div>
        </div>;
    }
}