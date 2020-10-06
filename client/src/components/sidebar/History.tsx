import * as React from "react";

import { classNames } from "../../utils";
import { ConnectionHandler } from "../../connection/ws";
import { GameHistoryRequest } from "../../requests/gameHistory";
import { EventState } from "../../objects/State";
import { Game } from "../../objects/Game";
import { userStorage } from "../../utils/storage";
import { ChipRequest } from "./history/ChipRequest";
import { GameHand } from "./history/GameHand";


type HistoryProps = {
    es: EventState;
    active: boolean;
    playerID: string;
    game?: Game;
    conn: ConnectionHandler;
}

type HistoryState = {
    history: any[];
}

export class History extends React.Component<HistoryProps, HistoryState> {
    constructor(props: any) {
        super(props)
        this.state = { history: [] }
    }

    public async componentDidMount() {
        let req = new GameHistoryRequest<any[]>();
        let history = await req.request({
            id: this.props.game.id.toString(),
            userToken: userStorage.getToken(),
        });

        console.log('history:', history);
        if (history) {
            this.setState({history});
        }
    }

    public render() {
        return <div className={classNames({"hidden": !this.props.active})}>
            <h3>Game History</h3>
            <div>
                {this.state.history.map((history) => {
                    let historyPlayer = this.props.es.getPlayer(history.player_id)

                    return history.hasOwnProperty("status") ? (
                        <ChipRequest {...history} player={this.props.es.getPlayer(this.props.playerID)} />
                    ) : (
                        <GameHand {...history} player={historyPlayer} />
                    );
                })}
            </div>
        </div>;
    }
}