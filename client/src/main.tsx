import * as React from "react";
import * as ReactDOM from "react-dom";

import { GameModal } from "./components/modals/Game";
import { userStorage } from "./utils/storage";
import { handleAuth, QPoker } from "./shared/entry";

let startGameButton = document.querySelector("#start-game-button")
let gameModal: GameModal;

let userStartGameEvent = () => {
    let userID = userStorage.getID();
    let userToken = userStorage.getToken();

    if (!userID || !userToken) {
        userStorage.removePlayer();
        QPoker.InitLogin();
        QPoker.OnPlayerFound.push(userStartGameEvent);
        return
    }

    QPoker.OnPlayerFound = QPoker.OnPlayerFound.slice(0, 1);
    gameModal.setState({isOpen: true});
}

startGameButton.addEventListener('click', userStartGameEvent);

ReactDOM.render(
    <GameModal
        ref={(comp) => { gameModal = comp; }}
    />,
    document.getElementById("game-modal")
);

handleAuth();
