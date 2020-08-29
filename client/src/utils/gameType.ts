import { GameTypesRequest } from "../requests/getGameTypes";
import { GameType } from "../objects/Game";
import { userStorage } from "./storage";


export const getGameTypes = (): Promise<GameType[]> => {
    let req = new GameTypesRequest<GameType[]>();

    return req.request({userToken: userStorage.getToken()});
}

export const getGameType = (gameTypeID: number): Promise<GameType> => {
    return new Promise<GameType>((resolve, reject) => {
        getGameTypes()
            .then((gameTypes: GameType[]) => {
                for (let i=0; i < gameTypes.length; i++) {
                    if (gameTypes[i].id === gameTypeID) {
                        resolve(gameTypes[i]);
                    }
                }

                reject("GameType not found: " + gameTypeID);
            }, reject);
    });
}
