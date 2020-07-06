import { createVideoAction, ClientAction } from "../utils";

class Player {
    conn: RTCPeerConnection;

    constructor(
        public fromPlayerID: number,
        public toPlayerID: number,
        public onEventCreated: (event: ClientAction) => void,
    ) {
        this.createConnection()
    }

    protected onIceCandidate(event: any) {
        console.log(typeof this, `${this.fromPlayerID} - ${this.toPlayerID}, onicecandidate:`, event.candidate, this.conn.connectionState);
    }

    protected onStateChange(event: any) {
        console.log(`${this.fromPlayerID} - ${this.toPlayerID}, connectionstatechange:`, event, this.conn.connectionState);
    }

    protected onIceStateChange(event: any) {
        console.log(`${this.fromPlayerID} - ${this.toPlayerID}, icestatechange:`, event);
    }

    protected onTrackEvent(event: any) {
        console.log(`${this.fromPlayerID} - ${this.toPlayerID}, track:`, event);

        if (this.getVideoElement().srcObject !== event.streams[0]) {
            this.getVideoElement().srcObject = event.streams[0];
        }
    }

    protected onNegotiationNeededEvent(event: any) {
        console.log(`${this.fromPlayerID} - ${this.toPlayerID}, negotiationneeded:`, event, this.conn.connectionState);
    }

    protected createConnection() {
        const configuration: any = null;
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
}

class RemotePlayer extends Player {
    hasCandidate: boolean;
    offer: RTCSessionDescription;
    candidate: RTCIceCandidate;

    protected onIceCandidate(event: any) {
        if (!event.candidate) {
            return;
        }

        super.onIceCandidate(event);
        this.hasCandidate = true;
        this.onEventCreated(
            createVideoAction({
                type: 'candidate',
                player_id: this.fromPlayerID,
                to_player_id: this.toPlayerID,
                candidate: event.candidate
            }
        ))
    }

    public handleOffer(offer: RTCSessionDescription): Promise<RTCSessionDescription> {
        if (this.offer && offer.sdp == this.offer.sdp) {
            return;
        }

        this.offer = offer;

        return new Promise<RTCSessionDescription>((resolve, reject) => {
            console.log('received offer:', offer, this.conn.connectionState);
            this.conn.setRemoteDescription(offer);
            this.conn.createAnswer()
                .then((offer: RTCSessionDescription) => {
                    console.log('sending answer:', offer);
                    this.conn.setLocalDescription(offer);
                    this.onEventCreated(
                        createVideoAction({
                            type: 'answer',
                            player_id: this.fromPlayerID,
                            to_player_id: this.toPlayerID,
                            offer: offer
                        })
                    );
                    resolve(offer);
                }, reject);
        });
    }

    public handleCandidate(candidate: RTCIceCandidate) {
        if (this.candidate && candidate.candidate == this.candidate.candidate) {
            return;
        }

        this.candidate = candidate;

        console.log("adding ice candidate:", candidate);
        this.conn.addIceCandidate(candidate);
    }
}

class LocalPlayer extends Player {
    answer: RTCSessionDescription;

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
                            player_id: this.fromPlayerID,
                            to_player_id: this.toPlayerID,
                            offer: offer
                        }
                    ))
                    resolve(offer);
                }, reject);
            
        });

    }

    public handleAnswer(answer: RTCSessionDescription): Promise<RTCSessionDescription> {
        if (this.answer && answer.sdp == this.answer.sdp) {
            return;
        }

        this.answer = answer;

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
        this.mediaStreamPromise = navigator.mediaDevices.getUserMedia({audio: false, video: true});
        this.userPlayer = new UserPlayer(this.playerID, this.playerID, this.onEventCreated)
        this.mediaStreamPromise.then(this.userPlayer.setStream.bind(this.userPlayer), console.error);
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
            if (this.playerID == playerID) {
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
        let playerID: any;

        this.setPlayers(event.data);

        for (playerID in this.remote) {
            let remotePlayer = this.remote[playerID];
            let toPlayer = event.data[playerID][this.playerID];

            if (toPlayer && toPlayer.offer) {
                remotePlayer.handleOffer(toPlayer.offer);
            }

            if (toPlayer && toPlayer.candidate) {
                remotePlayer.handleCandidate(toPlayer.candidate);
            }
        }

        for (playerID in this.local) {
            let localPlayer = this.local[playerID];
            let toPlayer = event.data[playerID][this.playerID];

            if (toPlayer && toPlayer.answer) {
                localPlayer.handleAnswer(toPlayer.answer);
            }
        }
    }
}
