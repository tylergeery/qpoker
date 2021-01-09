import * as React from "react";
import * as ReactDOM from "react-dom";

import { Nav } from "./components/Nav";
import { LoginModal } from "./components/modals/Login";
import { RegistrationModal } from "./components/modals/Registration";
import { LoginRequest } from "./requests/login";
import { RegistrationRequest } from "./requests/registration";
import { userStorage } from "./utils/storage";
import { Player } from "./objects/Player";
import { QPoker } from "./shared/entry";

var registrationModal: RegistrationModal;
var loginModal: LoginModal;
var nav: Nav;

let completePlayerAuth = (player: Player) => {
    userStorage.setUser(player);
    nav.setState({ player });

    QPoker.OnPlayerFound.map((fn: Function) => {
        fn(player);
    });
}

async function submitLogin (loginData: object) {
    const req = new LoginRequest<Player>();
    const player = await req.request({data: loginData})

    if (!req.success) {
        loginModal.setState({errors: req.errors});
        return
    }

    completePlayerAuth(player);
    loginModal.closeModal();
}

async function submitRegistration (loginData: object) {
    const req = new RegistrationRequest<Player>();
    const player = await req.request({data: loginData})

    if (!req.success) {
        registrationModal.setState({errors: req.errors});
        return
    }

    completePlayerAuth(player);
    registrationModal.closeModal();
}

const onLoginClick = () => {
    registrationModal.closeModal();
    loginModal.openModal();
}

const onRegisterClick = () => {
    loginModal.closeModal();
    registrationModal.openModal();
}

ReactDOM.render(
    <RegistrationModal
        ref={(comp) => { registrationModal = comp; }}
        register={submitRegistration}
        onLoginClick={onLoginClick}
    />,
    document.getElementById("registration-modal")
);

ReactDOM.render(
    <LoginModal
        ref={(comp) => { loginModal = comp; }}
        login={submitLogin}
        onRegisterClick={onRegisterClick}
    />,
    document.getElementById("login-modal")
);

ReactDOM.render(
    <Nav
        ref={(comp) => { nav = comp; }}
        onLoginClicked={onLoginClick}
        onRegisterClicked={onRegisterClick}
    />,
    document.getElementById("nav-account")
);

QPoker.InitLogin = onLoginClick;
