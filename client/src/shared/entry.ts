import { getPlayer } from "../helpers/player";
import { Game } from "../objects/Game";
import { Player } from "../objects/Player";
import { userStorage } from "../utils/storage";
import { VideoChannel } from "../video";

if (!window.hasOwnProperty('QPoker')) {
    console.error("Could not find QPoker config");
    throw Error("QPoker config not found");
}

type renderMethod = (game: Game, player: Player) => void;

declare global {
    interface Window {
        QPoker: {
            game: Game,
            OnPlayerFound: renderMethod[],
            InitLogin: () => void,
            VideoChannel?: VideoChannel,
        };
    }
}

export const QPoker = window.QPoker;
let game: Game = QPoker.game;

export const tableRender = async (render: renderMethod) => {
    let player = await getPlayer();

    if (!player) {
        userStorage.removePlayer();
        QPoker.InitLogin();
        QPoker.OnPlayerFound.push(tableRender.bind(null, render));

        // TODO: render table silhouette
        return
    }

    QPoker.OnPlayerFound = [];
    render(game, player)
}

// To communicate between entrypoints.
export const handleAuth = () => {
    QPoker.OnPlayerFound = QPoker.OnPlayerFound || [];
    QPoker.OnPlayerFound.push(userStorage.setUser.bind(userStorage));
};
