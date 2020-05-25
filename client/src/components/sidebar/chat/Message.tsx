import * as React from "react";

import { GamePlayer } from "../../../objects/State";

type MessageProps = {
    player: GamePlayer;
    message: boolean;
    ts?: string;
}

export class Message extends React.Component<MessageProps, {}> {
    public render() {
        return <div className="row">
            <div className="col s3"><b>{this.props.player ? this.props.player.username : 'Unknown'}</b></div>
            <div className="col s9">{this.props.message}</div>
        </div>;
    }
}