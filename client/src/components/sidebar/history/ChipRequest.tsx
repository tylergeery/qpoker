import * as React from "react";

import { GamePlayer } from "../../../objects/State";

type ChipRequestProps = {
    player: GamePlayer;
    status: string;
    amount: number;
}

export class ChipRequest extends React.Component<ChipRequestProps, {}> {
    public render() {
        return <div className="row">
            <div className="col s12">
            <b>{this.props.player ? this.props.player.username : 'Unknown'}</b>
            <span> ({this.props.amount})</span>
            <span> - {this.props.status}</span>
            </div>
        </div>;
    }
}