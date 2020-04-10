import * as React from "react";
import * as ReactDOM from "react-dom";

import { Table } from "./components/Table";

if (!window.hasOwnProperty('QPoker')) {
    console.error("Could not find QPoker config");
    throw Error("Qpoker config not found");
}

declare global {
    interface Window { QPoker: any; }
}

window.QPoker = window.QPoker || {};
let gameConfig: any = window.QPoker;

ReactDOM.render(
    <Table gameID={gameConfig.gameID} playerID={gameConfig.playerID} playerToken={gameConfig.playerToken} />,
    document.getElementById("table")
);