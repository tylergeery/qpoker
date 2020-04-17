class GameOptions {
    bigBlind: number;
    capacity: number;
}

export class Game {
    id: number;
    name: string;
    slug: string;
    ownerID: number;
    options: GameOptions;
    createdAt: string;
    updatedAt: string;
}
