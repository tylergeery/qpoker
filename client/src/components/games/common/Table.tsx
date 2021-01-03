import * as React from "react";

import { NewConnectionHandler, ConnectionHandler, EventType } from "../../../connection/ws";
import { Game } from "../../../objects/Game";
import { QPoker } from "../../../shared/entry";
import { VideoChannel } from "../../../video";

export type TableProps = {
    game?: Game;
    playerID: number;
    playerToken: string;
};

export abstract class VideoTable<EventState, State> extends React.Component<TableProps, State> {
    conn: ConnectionHandler;
    videoChannel: VideoChannel

    constructor(props: any) {
        super(props);
        this.state = this.getUpdatedState();
        this.resetConnection();
    }

    protected abstract getUpdatedState(evtState?: EventState): State;
    protected abstract needsSeat(): boolean;
    protected abstract formatGameEvent(data: any): EventState;

    protected enableVideo() {
        if (!this.props.playerID || this.videoChannel) {
            return
        }

        // TODO: check if video preferences have been turned on
        if (!this.needsSeat()) {
            this.videoChannel = new VideoChannel(
                +this.props.playerID,
                this.conn.send.bind(this.conn)
            );
            this.conn.subscribe(EventType.video, this.videoChannel.videoEvent.bind(this.videoChannel));

            // DEBUG
            QPoker.VideoChannel = this.videoChannel;
        }
    }

    protected stateUpdateHandler(evtState: EventState) {
        this.setState(this.getUpdatedState(evtState));
    }

    public resetConnection() {
        let initMsg = {
            token: this.props.playerToken,
            game_id: this.props.game.id
        };

        this.conn = NewConnectionHandler(this.formatGameEvent);
        this.conn.subscribe(EventType.game, this.stateUpdateHandler.bind(this));
        this.conn.init();
        this.conn.send(JSON.stringify(initMsg));
        this.conn.onDisconnect = () => {
            // wait for animation frame to not destroy background tabs
            window.requestAnimationFrame(this.resetConnection.bind(this));
        };
    }
};
