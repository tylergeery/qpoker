export class Game {
    id: number;
    name: string;
    slug: string;
    owner_id: number;
    game_type_id: number;
    options: {
        [key: string]: any;
    };
    created_at: string;
    updated_at: string;
}

export type Option = {
    game_type_id: number;
    game_option_id: number;
    name: string;
    label: string;
    type: string;
    default_value: any;
};

export type GameType = {
    id: number;
    key: string;
    display_name: string;
    options: Option[];
}
