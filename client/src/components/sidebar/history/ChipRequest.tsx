import * as React from "react";

import { AnonymousPlayer } from "../../../objects/Player";

type ChipRequestProps = {
    player: AnonymousPlayer;
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