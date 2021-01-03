export enum Suit {
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
