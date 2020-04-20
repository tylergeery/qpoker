import * as React from "react";
import * as ReactDOM from "react-dom";

import { Table } from "./components/Table";
import { userStorage } from "./utils/storage";
import { Game } from "./objects/Game";

if (!window.hasOwnProperty('QPoker')) {
    console.error("Could not find QPoker config");
    throw Error("Qpoker config not found");
}

declare global {
    interface Window { QPoker: any; }
}

window.QPoker = window.QPoker || {};
let game: Game = window.QPoker.game;

let tableRender = () => {
    let userID = userStorage.getID();
    let userToken = userStorage.getToken();

    if (!userID || !userToken) {
        userStorage.removePlayer();
        window.QPoker.InitLogin();
        window.QPoker.OnPlayerFound.push(tableRender);

        // TODO: render table silhouette
        return
    }

    window.QPoker.OnPlayerFound = [];
    ReactDOM.render(
        <Table game={game} playerID={userStorage.getID()} playerToken={userStorage.getToken()} />,
        document.getElementById("table")
    );
}

tableRender();