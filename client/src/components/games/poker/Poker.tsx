import * as React from "react";

import { classNames } from "../../../utils";
import { EventState, defaultEventState } from "./objects/State";
import { Player } from "./table/Player";
import { Seat } from "../common/Seat";
import { SideBar } from "../../SideBar";
import { ManageButtonSettings } from "../../sidebar/Settings";
import { VideoTable } from "../common/Table";
import { Chip } from "../common/Chip";
import { getCapacity } from "../../../objects/Player";


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

    protected getManageButtonSettings(): ManageButtonSettings {
        switch (this.state.es.manager.status) {
            case 'ready':
                return {text: 'start', disabled: false};
            case 'paused':
                return {text: 'resume', disabled: false};
            case 'active':
                return {text: 'pause', disabled: false};
            default:
                return {text: 'start', disabled: true};
        }
    }

    protected isPaused(): boolean {
        return this.state.es.manager.status === "paused";
    }

    protected getOffset(capacity: number): number {
        let i = 0, players = this.state.es.manager.state.table.players;

        for (; i < players.length; i++) {
            if (players[i] && players[i].id === this.props.playerID) {
                // put user at bottom of screen
                return capacity + (capacity/2) - i;
            }
        }

        return 0;
    }

    public render() {
        if (this.state.es.manager.state.table) {
            this.enableVideo();
        }

        const capacity = getCapacity(this.state.es.manager.state.table.players);
        const offset = this.getOffset(capacity);

        return (
            <div className="row white-text nmb">
                <div className="col s12 l9">
                    <div className={classNames(`row w100 table-holder nmb table-capacity-${capacity}`, {"paused": this.isPaused()})}>
                        <div className="board-holder">
                            <div>
                                {this.state.es.manager.state.board.map((card, i) =>
                                    <img key={i} className="card" src={`/assets/media/cards/${card.imageName()}.svg`} />
                                )}
                            </div>
                            <Chip amount={this.state.es.manager.pot.total} color="white" />
                        </div>
                        {this.state.es.manager.state.table ? this.state.es.manager.state.table.players.map((player: any, i: number) => {
                            return player ? 
                                <Player conn={this.conn}
                                        player={player}
                                        playerID={this.props.playerID}
                                        key={i}
                                        index={(i + offset) % capacity}
                                        manager={this.state.es.manager}
                                        game={this.props.game}
                                        cards={this.state.es.getPlayerCards(player.id)} />
                                : (
                                    (i < capacity) ? <Seat key={i} index={(i + offset) % capacity} /> : ''
                                );
                        }) : ''}
                        <img className="w100 bg" src="/assets/media/card_table.png" alt="Card table"/>
                    </div>
                </div>
                <div className="col s12 l3 sidebar-holder">
                    <SideBar {...this.props} conn={this.conn}
                        shouldRefreshHistory={this.state.es.refreshHistory}
                        manageButtonSettings={this.getManageButtonSettings()}
                        players={this.state.es.manager.state.table.players} />
                </div>
            </div>
        );
    }
}