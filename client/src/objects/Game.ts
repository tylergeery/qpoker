export class Game {
    id: number;
    name: string;
    slug: string;
    owner_id: number;
    options: {
        [key: string]: any;
    };
    created_at: string;
    updated_at: string;
}
