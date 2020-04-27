import * as React from "react";

import { classNames } from "../../utils";
import { ConnectionHandler } from "../../connection/ws";

type HistoryProps = {
    active: boolean;
    playerID: string;
    conn: ConnectionHandler;
}

type HistoryState = {
    history: any[];
}

export class History extends React.Component<HistoryProps, HistoryState> {
    constructor(props: any) {
        super(props)
        this.state = { history: [] }

        // TODO: register for new chats
    }

    public render() {
        return <div className={classNames({"hidden": !this.props.active})}>
            <h3>Game History</h3>
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