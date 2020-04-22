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
    id: number;
    stack: number;
    options: PlayerOptions;
    state: string;
    bigBlind: boolean;
    littleBlind: boolean;

    public static FromObj(playerObj: any): GamePlayer {
        let player = new GamePlayer();

        player.id = playerObj.id;
        player.stack = playerObj.stack;
        player.options = playerObj.options;
        player.state = playerObj.state;
        player.bigBlind = playerObj.big_blind;
        player.littleBlind = playerObj.little_blind;

        return player;
    }
}

class Player {
    id: number;
    username: string;
}
 
class Table {
    players: GamePlayer[];
    capacity: number;

    public static FromObj(tableObj: any): Table {
        let table = new Table();

        table.players = tableObj.players;
        table.capacity = tableObj.capacity;

        return table;
    }
}

class State {
    board: Card[];
    table: Table;
    state: string;

    public static FromObj(stateObj: any): State {
        let state = new State();

        state.board = stateObj.board.map((card: any) => {
            return new Card(
                card.value.toString(),
                String.fromCharCode(card.suit),
                String.fromCharCode(card.char),
            )
        });
        state.table = Table.FromObj(stateObj.table);
        state.state = stateObj.state;

        return state;
    }
}

type BetMap = {
    [key: number]: number;
}

class Pot {
    payouts: BetMap;
    playerBets: BetMap;
    playerTotals: BetMap;
    total: number;

    public static FromObj(potObj: any): Pot {
        let pot = new Pot();

        pot.payouts = potObj.payouts;
        pot.playerBets = potObj.player_bets;
        pot.playerTotals = potObj.player_totals;
        pot.total = potObj.total;

        return pot;
    }
}

class Manager {
    bigBlind: number;
    state: State;
    pot: Pot;

    public static FromObj(managerObj: any): Manager {
        let manager = new Manager();

        manager.bigBlind = managerObj.big_blind;
        manager.state = State.FromObj(managerObj.state);
        manager.pot = Pot.FromObj(managerObj.pot);

        return manager;
    }
}

export class EventState {
    manager: Manager;
    cards: {
        [key: number]: Card[];
    }
    players: {
        [key: number]: Player;
    }
    constructor(manager: Manager, cards: any, players: any) {
        this.manager = manager;
        this.cards = cards;
        this.players = players;
    }

    public static FromJSON(jsonStr: string): EventState {
        let {manager, cards, players} = JSON.parse(jsonStr);

        return new EventState(manager, cards, players);
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
        for (let i=0; i < this.manager.state.table.players.length; i++) {
            const player = this.manager.state.table.players[i];

            if (!player) {
                continue;
            }

            if (player.id.toString() !== playerID) {
                continue;
            }

            return player;
        }
    }
}

export const defaultEventState = EventState.FromJSON(`
{
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
        "big_blind": 0
    },
    "players": {},
    "cards": {}
}
`);
