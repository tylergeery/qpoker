export enum EventType {
    admin = 'admin',
    game = 'game',
    message = 'message',
    video = 'video',
}

type GameEventHandler = (data: any) => any;
type EventHandler = (action: any) => void;
type Subscribers = Map<EventType, EventHandler[]>;

export class ConnectionHandler {
    active: boolean;
    conn: WebSocket;
    subscribers: Subscribers;
    onDisconnect?: () => void;
    queue: string[];
    streak: number;
    gameEventHandler: GameEventHandler;

    constructor(conn: WebSocket, gameEventHandler: GameEventHandler) {
        this.conn = conn;
        this.active = false;
        this.queue = [];
        this.streak = 0;
        this.subscribers = this.initSubscribers();
        this.gameEventHandler = gameEventHandler;
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

    public subscribe(type: EventType, fn: (action: any) => void) {
        this.subscribers.get(type).push(fn);
    }

    private initSubscribers(): Subscribers {
        let subscribers: Subscribers = new Map();

        for (const eventType of Object.keys(EventType)) {
            subscribers.set(eventType as EventType, []);
        }

        return subscribers;
    }

    private publish(type: EventType, msg: any) {
        this.subscribers.get(type).map((fn) => fn(msg));
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
            case EventType.game:
                let state = this.gameEventHandler(event.data);

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

export function NewConnectionHandler(gameEventHandler: GameEventHandler): ConnectionHandler {
    const ws = new WebSocket(getHost());

    return new ConnectionHandler(ws, gameEventHandler);
}
