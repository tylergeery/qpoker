export type AnonymousPlayer = {
    id: number;
    username: string;
    state: string;
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

export const getCapacity = (players: (AnonymousPlayer | null)[]): number => {
    const allowedCapacities = [2, 4, 8, 12];
    const capacity = players.reduce((acc: number, player: AnonymousPlayer | null) => {
        if (!player) {
            return acc;
        }

        return acc + 1;
    }, 0);

    for (let i = 0; i < allowedCapacities.length; i++) {
        if (capacity === allowedCapacities[i]) {
            return capacity;
        }

        console.log("capacity:", capacity, allowedCapacities[i]);
        if (capacity < allowedCapacities[i]) {
            return allowedCapacities[i];
        }
    }

    return 12;
}
