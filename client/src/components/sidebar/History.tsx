import * as React from "react";

import { classNames } from "../../utils";
import { ConnectionHandler } from "../../connection/ws";
import { GameHistoryRequest } from "../../requests/gameHistory";
import { Game } from "../../objects/Game";
import { userStorage } from "../../utils/storage";
import { ChipRequest } from "./history/ChipRequest";
import { GameHand } from "./history/GameHand";
import { AnonymousPlayer, findPlayer } from "../../objects/Player";


type HistoryProps = {
    active: boolean;
    game?: Game;
    players: AnonymousPlayer[]
    playerID: number;
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
                {this.state.history.map((history, i) => {
                    let historyPlayer = findPlayer(history.player_id, this.props.players)

                    return history.hasOwnProperty("status") ? (
                        <ChipRequest key={i} {...history} player={findPlayer(this.props.playerID, this.props.players)} />
                    ) : (
                        <GameHand key={i} {...history} player={historyPlayer} />
                    );
                })}
            </div>
        </div>;
    }
}