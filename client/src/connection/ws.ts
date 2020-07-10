import { EventState } from "../objects/State";

export class ConnectionHandler {
    active: boolean;
    conn: WebSocket;
    subscribers: {
        [key: string]: ((action: any) => void)[];
    }
    onDisconnect?: () => void;
    queue: string[]
    streak: number;

    constructor(conn: WebSocket) {
        this.conn = conn;
        this.active = false;
        this.queue = [];
        this.streak = 0;
        this.subscribers = {
            'admin': [],
            'message': [],
            'game': [],
            'video': [],
        }
    }

    public init() {
        this.conn.onopen = (evt: Event) => {
            this.active = true;
            this.streak = 0;
            this.sendQueue();
        };

        this.conn.onerror = (evt: ErrorEvent) => {
            console.log("Connection error event: ", evt);
        };

        this.conn.onclose = (evt: CloseEvent) => {
            console.log("Connection close event: ", evt);
            this.active = false;
            this.streak++;

            if (this.streak <= 10 && this.onDisconnect) {
                this.onDisconnect();
            }
        };
    
        this.conn.onmessage = this.handleMessage.bind(this);
    }

    public subscribe(type: string, fn: (action: any) => void) {
        this.subscribers[type].push(fn);
    }

    private publish(type: string, msg: any) {
        this.subscribers[type].map((fn) => fn(msg));
    }

    public send(msg: any) {
        if (typeof msg != 'string') {
            msg = JSON.stringify(msg);
        }

        if (this.active) {
            this.conn.send(msg);
            return
        }

        this.queue.push(msg);
    }

    private handleMessage(message: MessageEvent) {
        let event = JSON.parse(message.data);

        switch(event.type) {
            case 'game':
                let state = EventState.FromObj(event.data);

                this.publish(event.type, state);
                break;
            default:
                this.publish(event.type, event)
        }
    }

    private sendQueue() {
        while (this.queue.length > 0) {
            let msg = this.queue.pop();

            this.conn.send(msg);
        }
    }
}

const getHost = (): string => {
    switch (window.location.hostname) {
        case 'localhost':
            return 'ws://localhost:8080/ws';
        default:
            return `wss://${window.location.hostname}/ws`;
    }
}

export function NewConnectionHandler(): ConnectionHandler {
    const ws = new WebSocket(getHost());

    return new ConnectionHandler(ws);
}
