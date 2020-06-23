enum Suit {
    BLANK = 'B',
    CLUBS = 'C',
    DIAMONDS = 'D',
    HEARTS = 'H',
    SPADES = 'S',
}

export class Card {
    constructor(public value: string, public suit: string, public char: string) {}

    public imageName(): string {
        return `${this.char}${this.suit}`;
    }
}

type PlayerOptions = {
    [key: string]: boolean;
}

export class GamePlayer {
    constructor(
        public id: number,
        public username: string,
        public stack: number,
        public options: PlayerOptions,
        public state: string,
        public bigBlind: boolean,
        public littleBlind: boolean
    ) {}

    public static FromObj(playerObj: any): GamePlayer {
        return new GamePlayer(
            playerObj.id,
            playerObj.username,
            playerObj.stack,
            playerObj.options,
            playerObj.state,
            playerObj.big_blind,
            playerObj.little_blind
        );
    }
}

export class Table {
    constructor(
        public players: GamePlayer[],
        public capacity: number,
        public active: number,
        public activeAt: number,
    ) {}

    public static FromObj(tableObj: any): Table {
        return new Table(
            tableObj.players, tableObj.capacity,
            tableObj.active, tableObj.active_at,
        );
    }
}

class State {
    constructor(
        public board: Card[],
        public table: Table,
        public state: string
    ) {}

    public static FromObj(stateObj: any): State {
        let state = new State(
            (stateObj.board || []).map((card: any) => {
                return new Card(
                    card.value.toString(),
                    String.fromCharCode(card.suit),
                    String.fromCharCode(card.char),
                )
            }),
            Table.FromObj(stateObj.table),
            stateObj.state
        );

        return state;
    }
}

type BetMap = {
    [key: number]: number;
}

class Pot {
    constructor(
        public payouts: BetMap,
        public playerBets: BetMap,
        public playerTotals: BetMap,
        public total: number
    ) {}

    public static FromObj(potObj: any): Pot {
        return new Pot(
            potObj.payouts || {},
            potObj.player_bets || {},
            potObj.player_totals || {},
            potObj.total
        );
    }

    public toCover(playerID: number): number {
        return this.maxBet() - this.playerBets[playerID];
    }

    private maxBet(): number {
        let playerID: any = 0;
        let maxBet = 0;

        for (playerID in this.playerBets) {
            maxBet = Math.max(maxBet, this.playerBets[playerID]);
        }

        return maxBet;
    }
}

export class Manager {
    constructor(
        public bigBlind: number,
        public state: State,
        public pot: Pot,
        public status: string,
        public allIn: boolean
    ) {}

    public static FromObj(managerObj: any): Manager {
        return new Manager(
            managerObj.big_blind,
            State.FromObj(managerObj.state),
            Pot.FromObj(managerObj.pot || {}),
            managerObj.status,
            managerObj.all_in,
        );
    }
}

export class EventState {
    constructor(
        public manager: Manager,
        public cards: {[key: number]: Card[]},
    ) {}

    public static FromObj(obj: any) {
        for (let playerID in obj.cards) {
            obj.cards[playerID] = obj.cards[playerID].map((card: any) => {
                return new Card(
                    card.value.toString(),
                    String.fromCharCode(card.suit),
                    String.fromCharCode(card.char),
                )
            })
        }

        return new EventState(
            Manager.FromObj(obj.manager),
            obj.cards);
    }

    public getPlayerCards(playerID: number): Card[] {
        if (!this.cards[playerID]) {
            return [
                new Card('1', Suit.BLANK, '1'),
                new Card('1', Suit.BLANK, '1')
            ];
        }

        return this.cards[playerID];
    }

    public getPlayer(playerID: string): GamePlayer {
        if (!playerID) {
            return null
        }

        for (let i=0; i < this.manager.state.table.players.length; i++) {
            const player = this.manager.state.table.players[i];

            if (!player) {
                continue;
            }

            if (player.id.toString() !== playerID.toString()) {
                continue;
            }

            return player;
        }

        return null;
    }
}

export const defaultEventState = EventState.FromObj({
    "manager": {
        "game_id": 0,
        "state": {
            "board": [],
            "table": {
                "players": [],
                "capacity": 0
            },
            "state": "init"
        },
        "pot": {
            "payouts": {},
            "player_bets": {},
            "player_totals": {},
            "total": 0
        },
        "big_blind": 0,
        "status": "init",
        "all_in": false,
    },
    "cards": {}
});
