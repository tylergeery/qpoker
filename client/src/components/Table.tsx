import * as React from "react";

import { EventState, defaultEventState } from "../objects/State";
import { NewConnectionHandler, ConnectionHandler } from "../connection/ws";
import { Game } from "../objects/Game";
import { Player } from "./Player";
import { TaskBar } from "./TaskBar";
import { SideBar } from "./SideBar";

export type TableProps = {
    game?: Game;
    playerID: string;
    playerToken: string;
}

export type TableState = {
    es: EventState;
}
// State is never set so we use the '{}' type.
export class Table extends React.Component<TableProps, TableState> {
    conn: ConnectionHandler;

    constructor(props: any) {
        super(props);
        this.state = {es: defaultEventState};
    }

    public componentDidMount() {
        this.resetConnection();
    }

    public resetConnection() {
        let initMsg = {
            token: this.props.playerToken,
            game_id: this.props.game.id
        };

        this.conn = NewConnectionHandler(this.stateUpdateHandler.bind(this));
        this.conn.send(JSON.stringify(initMsg));
        this.conn.onDisconnect = this.resetConnection.bind(this);
    }

    public sendAction(action: any) {
        this.conn.send(JSON.stringify(action));
    }

    public stateUpdateHandler(evt: MessageEvent) {
        this.setState({es: EventState.FromJSON(evt.data)});
    }

    public render() {
        return (
            <div className="row grey lighten-4 white-text nmb">
                <div className="col s12 m7 l9">
                    <div className="row w100 table-holder nmb">
                        <h5>{this.props.game ? this.props.game.name : ''}</h5>
                        {this.state.es.manager.state.table ? this.state.es.manager.state.table.players.map((player: any, i: number) => {
                            return player ? <Player player={player} index={i} gameState={this.state.es.manager.state.state} cards={this.state.es.getPlayerCards(player.id)} /> : '';
                        }) : ''}
                        <img className="w100 bg" src="/assets/media/card_table.png" alt="Card table"/>
                        <TaskBar
                            player={this.state.es.getPlayer(this.props.playerID)}
                            game={this.props.game}
                            sendAction={this.sendAction.bind(this)} />
                    </div>
                </div>
                <div className="col s12 m5 l3 sidebar-holder">
                        <SideBar
                            sendAction={this.sendAction.bind(this)}
                            {...this.props} />
                </div>
            </div>
        );
    }
}