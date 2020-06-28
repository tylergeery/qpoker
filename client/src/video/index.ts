class Player {
    id: string;
    video: any;

    constructor(playerID: string) {
        this.id = playerID;
        this.video = this.getVideoElement()
    }

    protected getVideoElement(): HTMLVideoElement {
        return document.querySelector(`#player-video-${this.id}`);
    }
}

class LocalPlayer extends Player {
    conn: RTCPeerConnection;

    constructor(playerID: string) {
        super(playerID)
        this.getMedia();
        this.createConnection();
    }

    protected onIceCandidate(event: any) {
        // send ice candidate to others
    }

    protected onIceStateChange(event: any) {
        console.log('icestatechange:', event);
    }

    protected onTrackEvent(event: any) {
        console.log('track:', event);
    }

    private getMedia() {
        navigator.mediaDevices.getUserMedia({audio: true, video: true})
            .then((stream) => {
                console.log('gotmediastream:', stream);
                this.getVideoElement().srcObject = stream;
                stream.getTracks().forEach(track => this.conn.addTrack(track, stream));
            }, console.error);
    }

    private createConnection() {
        this.conn = new RTCPeerConnection({});
        this.conn.addEventListener('icecandidate', this.onIceCandidate.bind(this));
        this.conn.addEventListener('iceconnectionstatechange', this.onIceStateChange.bind(this));
        this.conn.addEventListener('track', this.onTrackEvent.bind(this));
    }

    public startVideo(): Promise<RTCSessionDescription> {
        return new Promise<RTCSessionDescription>((resolve, reject) => {
            this.conn.createOffer()
                .then((offer: RTCSessionDescription) => {
                    this.conn.setLocalDescription(offer);
                    resolve(offer);
                }, reject);
        });
    }
}

export class VideoChannel {
    player: LocalPlayer;
    players: {
        [playerID: string]: Player
    };

    constructor(playerID: string) {
        this.player = new LocalPlayer(playerID);
        this.players = {}; // TODO
    }

    public startVideo(): Promise<RTCSessionDescription> {
        return this.player.startVideo();
    }
}
