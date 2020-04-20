class GameOptions {
    big_blind: number;
    capacity: number;
}

export class Game {
    id: number;
    name: string;
    slug: string;
    owner_id: number;
    options: GameOptions;
    created_at: string;
    updated_at: string;
}
    