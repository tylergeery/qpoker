class GameOptions {
    public big_blind: number;
    public capacity: number;
    public decision_time: number;
    public time_between_hands: number;
    public buy_in_max: number;
    public buy_in_min: number;
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
