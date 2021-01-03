import * as React from "react";
import { AnonymousPlayer } from "../../../objects/Player";

type GameHandProps = {
    player: AnonymousPlayer;
    board: any;
    cards: any;
    bets: any
    ending: number;
    starting: number;
}

export class GameHand extends React.Component<GameHandProps, {}> {
    public render() {
        return <div className="row">
            <div className="col s12">
                <b>{this.props.player ? this.props.player.username : 'Unknown'}</b>
                <span> ({this.props.cards})</span>
            </div>
        </div>;
    }
}