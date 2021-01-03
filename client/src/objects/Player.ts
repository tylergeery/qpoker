export type AnonymousPlayer = {
    id: number;
    username: string;
};

export type PlayerWithChips = {
    stack: number;
};

export type Player = {
    id: number;
    email: string;
    token: string;
    username: string;
    created_at: string;
    updated_at: string;
};

export type AnonymousPlayerWithChips = AnonymousPlayer & PlayerWithChips;

export const findPlayer = <PlayerType extends AnonymousPlayer>(
    playerID: number,
    players: PlayerType[],
): PlayerType | null => {
    if (!playerID) {
        return null
    }

    for (const player of players) {
        if (!player) {
            continue;
        }

        if (player.id.toString() !== playerID.toString()) {
            continue;
        }

        return player;
    }

    return null;
};
