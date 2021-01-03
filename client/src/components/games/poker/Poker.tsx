import * as React from "react";

import { EventState, defaultEventState } from "./objects/State";
import { Player } from "./table/Player";
import { Seat } from "../common/Seat";
import { SideBar } from "../../SideBar";
import { VideoTable } from "../common/Table";


export type TableState = {
    es: EventState;
}

// State is never set so we use the '{}' type.
export class Table extends VideoTable<EventState, TableState> {

    protected getUpdatedState(evtState?: EventState): TableState {
        return { es: evtState ?? defaultEventState };
    }

    protected formatGameEvent(data: any): EventState {
        return EventState.FromObj(data)
    }
    protected needsSeat(): boolean {
        let i = 0, players = this.state.es.manager.state.table.players;

        for (; i < players.length; i++) {
            if (!players[i]) {
                continue;
            }

            if (players[i].id.toString() === this.props.playerID.toString()) {
                return false;
            }
        }

        return true
    }

    public render() {
        if (this.state.es.manager.state.table) {
            this.enableVideo();
        }

        return (
            <div className="row white-text nmb">
                <div className="col s12 l9">
                    <div className="row w100 table-holder nmb">
                        <div className="board-holder">
                            {this.state.es.manager.state.board.map((card, i) =>
                                <img key={i} className="card" src={`/assets/media/cards/${card.imageName()}.svg`} />
                            )}
                        </div>
                        {this.state.es.manager.state.table ? this.state.es.manager.state.table.players.map((player: any, i: number) => {
                            return player ? 
                                <Player conn={this.conn}
                                        player={player}
                                        playerID={this.props.playerID}
                                        key={i}
                                        index={i}
                                        manager={this.state.es.manager}
                                        game={this.props.game}
                                        cards={this.state.es.getPlayerCards(player.id)} />
                                : <Seat key={i} index={i} />;
                        }) : ''}
                        <img className="w100 bg" src="/assets/media/card_table.png" alt="Card table"/>
                    </div>
                </div>
                <div className="col s12 l3 sidebar-holder">
                    <SideBar {...this.props} conn={this.conn}
                        showStartButton={this.state.es.manager.status != "active"}
                        disableStartButton={this.state.es.manager.status == "init"}
                        players={this.state.es.manager.state.table.players} />
                </div>
            </div>
        );
    }
}