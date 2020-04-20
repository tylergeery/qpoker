import * as React from "react";

import { EventState, defaultEventState } from "../objects/State";
import { NewConnectionHandler, ConnectionHandler } from "../connection/ws";
import { Game } from "../objects/Game";

export type TableProps = {
    game?: Game;
    playerID: string;
    playerToken: string;
}

// State is never set so we use the '{}' type.
export class Table extends React.Component<TableProps, EventState> {
    conn: ConnectionHandler;

    constructor(props: any) {
        super(props);
        this.state = defaultEventState;
    }

    public componentDidMount() {
        let initMsg = {
            token: this.props.playerToken,
            game_id: this.props.game.id
        };

        this.conn = NewConnectionHandler(this.stateUpdateHandler.bind(this));
        this.conn.send(JSON.stringify(initMsg));
    }

    public stateUpdateHandler(evt: MessageEvent) {
        this.setState(EventState.FromJSON(evt.data));
    }

    public render() {
        return (
            <div>
                {this.state.manager.state.table ? this.state.manager.state.table.players.length : 0}
                <h4>Table big blind is: {this.props.game.options.big_blind} and capacity: {this.props.game.options.capacity}!</h4>
                <div className="row w100 table-holder">
                    {this.state.manager.state.table ? this.state.manager.state.table.players.map((player: any, i: number) => {
                        return <div className={ `player player-${i}` }>
                            {player ? 'Player ' + player.id : ''}
                            <img className="card" src="/assets/media/cards/1B.svg" />
                            <img className="card" src="/assets/media/cards/1J.svg" />
                        </div>;
                    }) : ''}
                    <img className="w100 bg" src="/assets/media/card_table.png" alt="Card table"/>
                </div>
            </div>
        );
    }
}