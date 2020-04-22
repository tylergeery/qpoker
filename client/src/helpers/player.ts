import { userStorage } from "../utils/storage";
import { Player } from "../objects/Player";
import { GetPlayerRequest } from "../requests/getPlayer";


export async function getPlayer(): Promise<Player> {
    let userID = userStorage.getID();
    let userToken = userStorage.getToken();

    if (!userID || !userToken) {
        return null;
    }

    let req = new GetPlayerRequest<Player>();
    let player = await req.request({id: userID, userToken});
    if (!req.success) {
        return null;
    }

    return player;
}