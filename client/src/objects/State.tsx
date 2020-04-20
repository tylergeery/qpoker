var suits = new String("CDHS")
enum Suit {
    CLUBS = suits.charCodeAt(0),
    DIAMONDS = suits.charCodeAt(1),
    SPADES = suits.charCodeAt(2),
    HEARTS = suits.charCodeAt(3),
}

class Card {
    constructor(public value: string, public suit: Suit) {}
}

type PlayerOptions = {
    [key: string]: boolean;
}

class GamePlayer {
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

        state.board = stateObj.board;
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

        console.log("json:", jsonStr, manager, cards, players);
        return new EventState(manager, cards, players);
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
