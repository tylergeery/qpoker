import * as React from "react";
import * as ReactDOM from "react-dom";

import { GameModal } from "./components/modals/Game";
import { userStorage } from "./utils/storage";

let startGameButton = document.querySelector("#start-game-button")
let gameModal: GameModal;

let userStartGameEvent = () => {
    let userID = userStorage.getID();
    let userToken = userStorage.getToken();

    if (!userID || !userToken) {
        userStorage.removePlayer();
        window.QPoker.InitLogin();
        window.QPoker.OnPlayerFound.push(userStartGameEvent);
        return
    }

    window.QPoker.OnPlayerFound = window.QPoker.OnPlayerFound.slice(0, 1);
    gameModal.setState({isOpen: true});
}

startGameButton.addEventListener('click', userStartGameEvent);

ReactDOM.render(
    <GameModal
        ref={(comp) => { gameModal = comp; }}
    />,
    document.getElementById("game-modal")
);

// To communicate between entrypoints. TODO: move to a shared module
window.QPoker.OnPlayerFound = window.QPoker.OnPlayerFound || [];
window.QPoker.OnPlayerFound.push(userStorage.setUser.bind(userStorage));
