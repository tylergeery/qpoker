import * as React from "react";

import { AnonymousPlayer } from "../../../objects/Player";

type MessageProps = {
    playerUsername: string;
    message: string;
    ts?: string;
};

export class Message extends React.Component<MessageProps, {}> {
    public render() {
        return <div className="row">
            <div className="col s3"><b>{this.props.playerUsername}</b></div>
            <div className="col s9">{this.props.message}</div>
        </div>;
    }
}
