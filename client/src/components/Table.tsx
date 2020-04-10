import * as React from "react";

import { Manager, ManagerDefault, ManagerFromJSON } from "../objects/State";
import { NewConnectionHandler, ConnectionHandler } from "../connection/ws";

export type TableProps = {
    playerID: number;
    playerToken: string;
    gameID: string;
}

export type TableState = {
    manager: Manager;
}

// State is never set so we use the '{}' type.
export class Table extends React.Component<TableProps, TableState> {
    conn: ConnectionHandler;

    constructor(props: any) {
        super(props);
        this.state = { manager: ManagerDefault };
    }

    public componentDidMount() {
        this.conn = NewConnectionHandler(this.stateUpdateHandler);
        this.conn.send(JSON.stringify(this.props));
    }

    public stateUpdateHandler(evt: MessageEvent) {
        const nextManager: Manager = ManagerFromJSON(evt.data);

        this.setState((prevState: TableState) => {
            return { manager: nextManager };
        })
    }

    public render() {
        return <h1>Table big blind is: {this.state.manager.big_blind} and total pot: {this.state.manager.pot.total}!</h1>;
    }
}