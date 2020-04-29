import * as React from "react";

import { EventState, defaultEventState } from "../objects/State";
import { NewConnectionHandler, ConnectionHandler } from "../connection/ws";
import { Game } from "../objects/Game";
import { Player } from "./table/Player";
import { Seat } from "./table/Seat";
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
        this.resetConnection();
    }

    public resetConnection() {
        let initMsg = {
            token: this.props.playerToken,
            game_id: this.props.game.id
        };

        this.conn = NewConnectionHandler();
        this.conn.subscribe('game', this.stateUpdateHandler.bind(this));
        this.conn.init();
        this.conn.send(JSON.stringify(initMsg));
        this.conn.onDisconnect = this.resetConnection.bind(this);
    }

    public stateUpdateHandler(evtState: EventState) {
        this.setState({es: evtState});
    }

    public needsSeat(): boolean {
        let i = 0, players = this.state.es.manager.state.table.players;

        for (; i < players.length; i++) {
            if (!players[i]) {
                continue;
            }

            if (players[i].id.toString() === this.props.playerID) {
                return false;
            }
        }

        return true
    }

    public render() {
        let chooseSeat = this.needsSeat();

        return (
            <div className="row white-text nmb">
                <div className="col s12 l9">
                    <div className="row w100 table-holder nmb">
                        <div className="board-holder">
                            {this.state.es.manager.state.board.map((card) => {
                                return <img className="card" src={`/assets/media/cards/${card.imageName()}.svg`} />
                            })}
                        </div>
                        {this.state.es.manager.state.table ? this.state.es.manager.state.table.players.map((player: any, i: number) => {
                            return player ? 
                                <Player conn={this.conn} player={player} playerID={this.props.playerID} index={i} gameState={this.state.es.manager.state.state} cards={this.state.es.getPlayerCards(player.id)} />
                                : (chooseSeat ? <Seat index={i} /> : '');
                        }) : ''}
                        <img className="w100 bg" src="/assets/media/card_table.png" alt="Card table"/>
                    </div>
                </div>
                <div className="col s12 l3 sidebar-holder">
                    <SideBar {...this.props} conn={this.conn} es={this.state.es}/>
                </div>
            </div>
        );
    }
}