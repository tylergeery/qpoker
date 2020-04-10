var suits = new String("CDHS")
enum Suit {
    CLUBS = suits.charCodeAt(0),
    DIAMONDS = suits.charCodeAt(1),
    SPADES = suits.charCodeAt(2),
    HEARTS = suits.charCodeAt(3),
}

class Card {
    value: string;
    suit: Suit;   
}

type PlayerOptions = {
    [key: string]: boolean;
}

class Player {
    id: number;
    stack: number;
    options: PlayerOptions;
    state: string;
    big_blind: boolean;
    little_blind: boolean;
}

class Table {
    players: Player[];
    capacity: number;
    Active: number;
    Dealer: number;
}

class HoldEm {
    board: Card[];
    table: Table;
    state: string;
}

type BetMap = {
    [key: number]: number;
}

class Pot {
    payouts: BetMap;
    player_bets: BetMap;
    player_totals: BetMap;
    total: number;
}

export class Manager {
    big_blind: number;
    state: HoldEm;
    pot: Pot;
    constructor(big_blind: number, state: HoldEm, pot: Pot) {
        this.big_blind = big_blind;
        this.state = state;
        this.pot = pot;
    }
}

export function ManagerFromJSON(jsonStr: string): Manager {
    let manager: Manager = JSON.parse(jsonStr);

    //TODO: validate

    return manager;
}

const pot = new Pot();
const holdem = new HoldEm();
export const ManagerDefault = new Manager(100, holdem, pot);
