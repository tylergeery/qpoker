export class ConnectionHandler {
    active: boolean;
    conn: WebSocket;
    onMessageEvent: (evt: MessageEvent) => void;
    queue: string[]

    constructor(conn: WebSocket, onMessageEvent: (evt: MessageEvent) => void) {
        this.conn = conn;
        this.onMessageEvent = onMessageEvent;
        this.active = false;
        this.queue = [];

        this.setup();
    }

    private setup() {
        this.conn.onopen = (evt: Event) => {
            console.log("Connection open event: ", evt);
            this.sendQueue();
        };

        this.conn.onerror = (evt: ErrorEvent) => {
            console.log("Connection error event: ", evt);
        };

        this.conn.onclose = (evt: CloseEvent) => {
            console.log("Connection close event: ", evt);
        };
    
        this.conn.onmessage = (evt: MessageEvent) => {
            console.log("Connection message event: ", evt);
            this.onMessageEvent(evt);
        };
    }

    public send(msg: string) {
        if (this.active) {
            console.log("Sending message event:", msg)
            this.conn.send(msg);
            return
        }

        console.log("Queueing message event:", msg)
        this.queue.push(msg);
    }

    private sendQueue() {
        while (this.queue.length > 0) {
            let msg = this.queue.pop();

            console.log("Sending message event:", msg)
            this.conn.send(msg);
        }
    }
}

export function NewConnectionHandler(onMessageEvent: (evt: MessageEvent) => void): ConnectionHandler {
    const ws = new WebSocket('ws://localhost:8080/ws');

    return new ConnectionHandler(ws, onMessageEvent);
}
