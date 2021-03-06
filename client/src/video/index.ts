import { createVideoAction, ClientAction } from "../utils";

abstract class Player {
    conn: RTCPeerConnection;
    prefix: string;

    constructor(
        public fromPlayerID: number,
        public toPlayerID: number,
        public onEventCreated: (event: ClientAction) => void,
    ) {
        this.createConnection();
        this.prefix = this.getPrefix();
    }

    protected abstract getPrefix(): string;

    protected abstract getCandidateEventName(): string;

    protected onIceCandidate(event: any) {
        console.log(this.prefix, 'onicecandidate:', event.candidate, this.conn.connectionState);

        if (!event.candidate) {
            return;
        }

        this.onEventCreated(
            createVideoAction({
                type: this.getCandidateEventName(),
                from_player_id: this.fromPlayerID,
                to_player_id: this.toPlayerID,
                candidate: event.candidate
            }
        ));
    }

    protected onStateChange(event: any) {
        console.log(this.prefix, 'connectionstatechange:', event, this.conn.connectionState);
    }

    protected onIceStateChange(event: any) {
        console.log(this.prefix, 'icestatechange:', event);

        if (event.target.iceConnectionState != "failed") {
            return;
        }

        // TODO: try to reinitiate
    }

    protected onTrackEvent(event: any) {
        console.log(this.prefix, 'track:', event);

        if (this.getVideoElement().srcObject !== event.streams[0]) {
            this.getVideoElement().srcObject = event.streams[0];
        }
    }

    protected onNegotiationNeededEvent(event: any) {
        console.log(this.prefix, 'negotiationneeded:', event, this.conn.connectionState);
    }

    protected createConnection() {
        const configuration: any = {'iceServers': [{ 'urls': 'stun:stun.l.google.com:19302' }]};
        this.conn = new RTCPeerConnection(configuration);
        this.conn.onicecandidate = this.onIceCandidate.bind(this);
        this.conn.onconnectionstatechange = this.onStateChange.bind(this);
        this.conn.oniceconnectionstatechange = this.onIceStateChange.bind(this);
        this.conn.ontrack = this.onTrackEvent.bind(this);
        this.conn.onnegotiationneeded = this.onNegotiationNeededEvent.bind(this);
    }

    protected getVideoElement(): HTMLVideoElement {
        return document.querySelector(`#player-video-${this.toPlayerID}`);
    }

    public mute() {
        this.getVideoElement().muted = true;
    }

    public mirrorStream() {
        this.getVideoElement().style.transform = 'scale(-1, 1)';
    }

    public handleCandidate(candidate: RTCIceCandidate) {
        console.log(this.prefix, 'adding ice candidate:', candidate);
        this.conn.addIceCandidate(candidate)
            .then((candidate) => {
                console.log(this.prefix, 'added ice candidate:', candidate);
            }, console.error);
    }
}

class RemotePlayer extends Player {
    protected getPrefix(): string {
        return `${this.toPlayerID} - ${this.fromPlayerID}`;
    }

    protected getCandidateEventName(): string {
        return 'candidate-remote';
    }

    public handleOffer(offer: RTCSessionDescription): Promise<RTCSessionDescription> {
        return new Promise<RTCSessionDescription>((resolve, reject) => {
            console.log(this.prefix, 'received offer:', offer, this.conn.connectionState);
            this.conn.setRemoteDescription(offer);
            this.conn.createAnswer()
                .then((offer: RTCSessionDescription) => {
                    console.log('sending answer:', offer);
                    this.conn.setLocalDescription(offer);
                    this.onEventCreated(
                        createVideoAction({
                            type: 'answer',
                            from_player_id: this.fromPlayerID,
                            to_player_id: this.toPlayerID,
                            offer: offer
                        })
                    );
                    resolve(offer);
                }, reject);
        });
    }
}

class LocalPlayer extends Player {
    protected getPrefix(): string {
        return `${this.fromPlayerID} - ${this.toPlayerID}`;
    }

    protected getCandidateEventName(): string {
        return 'candidate-local';
    }

    public createOffer(stream: MediaStream): Promise<RTCSessionDescription> {
        stream.getTracks().forEach(track => this.conn.addTrack(track, stream));

        return new Promise<RTCSessionDescription>((resolve, reject) => {
            console.log("creating offer");
            this.conn.createOffer()
                .then((offer: RTCSessionDescription) => {
                    this.conn.setLocalDescription(offer);
                    this.onEventCreated(
                        createVideoAction({
                            type: 'offer',
                            from_player_id: this.fromPlayerID,
                            to_player_id: this.toPlayerID,
                            offer: offer
                        }
                    ))
                    resolve(offer);
                }, reject);
            
        });

    }

    public handleAnswer(answer: RTCSessionDescription): Promise<RTCSessionDescription> {
        return new Promise<RTCSessionDescription>((resolve, reject) => {
            console.log('received answer:', answer, this.conn.connectionState);
            this.conn.setRemoteDescription(answer);
        });
    }
}

class UserPlayer extends LocalPlayer {
    public setStream(stream: MediaStream) {
        this.getVideoElement().srcObject = stream;
    }

    protected onTrackEvent(event: any) {}
    protected onIceCandidate(event: any) {}
    protected onStateChange(event: any) {}
    protected onIceStateChange(event: any) {}
    protected onNegotiationNeededEvent(event: any) {}
}

export class VideoChannel {
    userPlayer: UserPlayer;
    local: {
        [playerID: number]: LocalPlayer
    };
    remote: {
        [playerID: number]: RemotePlayer
    };
    mediaStreamPromise: Promise<MediaStream>;

    constructor(
        public playerID: number,
        public onEventCreated: (event: ClientAction) => void
    ) {
        this.remote = {};
        this.local = {};
        this.mediaStreamPromise = navigator.mediaDevices.getUserMedia({audio: true, video: true});
        this.userPlayer = new UserPlayer(this.playerID, this.playerID, this.onEventCreated)
        this.mediaStreamPromise.then(
            (stream: MediaStream) => {
                this.userPlayer.setStream(stream);
                this.userPlayer.mute();
                this.userPlayer.mirrorStream();
            },
            console.error,
        );
    }

    public setPlayers(playerIDOffers: object) {
        let playerID: any;

        // Remove players who have left game
        for (playerID in this.remote) {
            if (!playerIDOffers.hasOwnProperty(playerID)) {
                delete this.remote[playerID];
            }
        }

        for (playerID in this.local) {
            if (!playerIDOffers.hasOwnProperty(playerID)) {
                delete this.local[playerID];
            }
        }

        // Create new players as needed
        for (playerID in playerIDOffers) {
            playerID = +playerID;
            if (!playerID || this.playerID == playerID) {
                continue
            }

            if (!this.remote[playerID]) {
                this.remote[playerID] = new RemotePlayer(this.playerID, playerID, this.onEventCreated);
            }

            if (!this.local[playerID]) {
                this.local[playerID] = new LocalPlayer(this.playerID, playerID, this.onEventCreated);
                this.mediaStreamPromise.then(
                    this.local[playerID].createOffer.bind(this.local[playerID]),
                    console.error
                );
            }
        }
    }

    public videoEvent(event: ClientAction) {
        let remotePlayer: RemotePlayer;
        let localPlayer: LocalPlayer;

        switch (event.data.type) {
            case 'offer':
                remotePlayer = this.remote[+event.data.from_player_id];
                remotePlayer.handleOffer(event.data.offer);
                break;
            case 'answer':
                localPlayer = this.local[+event.data.from_player_id];
                localPlayer.handleAnswer(event.data.offer);
                break;
            case 'candidate-local':
                remotePlayer = this.remote[+event.data.from_player_id];
                remotePlayer.handleCandidate(event.data.candidate);
                break;
            case 'candidate-remote':
                localPlayer = this.local[+event.data.from_player_id];
                localPlayer.handleCandidate(event.data.candidate);
                break;
            default:
                this.setPlayers(event.data);
                break;
        }
    }
}
