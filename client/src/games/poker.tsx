import * as React from "react";
import * as ReactDOM from "react-dom";
import { Table } from "../components/games/poker/Poker";
import { Game } from "../objects/Game";
import { Player } from "../objects/Player";
import { handleAuth, tableRender } from "../shared/entry";

const render = (game: Game, player: Player) => {
    ReactDOM.render(
        <Table game={game} playerID={+player.id} playerToken={player.token} />,
        document.getElementById("table"),
    );
};

tableRender(render);
handleAuth();
