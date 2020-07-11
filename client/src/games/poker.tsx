import * as React from "react";
import * as ReactDOM from "react-dom";

import { Table } from "../components/games/poker/Poker";
import { userStorage } from "../utils/storage";
import { Game } from "../objects/Game";
import { getPlayer } from "../helpers/player";

if (!window.hasOwnProperty('QPoker')) {
    console.error("Could not find QPoker config");
    throw Error("Qpoker config not found");
}

declare global {
    interface Window { QPoker: any; }
}

window.QPoker = window.QPoker || {};
let game: Game = window.QPoker.game;

let tableRender = async () => {
    let player = await getPlayer();

    if (!player) {
        userStorage.removePlayer();
        window.QPoker.InitLogin();
        window.QPoker.OnPlayerFound.push(tableRender);

        // TODO: render table silhouette
        return
    }

    window.QPoker.OnPlayerFound = [];
    ReactDOM.render(
        <Table game={game} playerID={player.id.toString()} playerToken={player.token} />,
        document.getElementById("table")
    );
}

tableRender();