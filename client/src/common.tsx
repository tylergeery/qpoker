import * as React from "react";
import * as ReactDOM from "react-dom";

import { Nav } from "./components/Nav";
import { LoginModal } from "./components/modals/Login";
import { RegistrationModal } from "./components/modals/Registration";
import { LoginRequest } from "./requests/login";
import { RegistrationRequest } from "./requests/registration";
import { userStorage } from "./utils/storage";
import { Player } from "./objects/Player";

var registrationModal: RegistrationModal;
var loginModal: LoginModal;
var nav: Nav;

let completePlayerAuth = (player: Player) => {
    userStorage.setUser(player);
    nav.setState({ player });

    window.QPoker.OnPlayerFound.map((fn: Function) => {
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

const onLoginClick = (event: any) => {
    registrationModal.closeModal();
    loginModal.openModal();
}

const onRegisterClick = (event: any) => {
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

// To communicate between entrypoints. TODO: move to a shared module for consistency
window.QPoker.OnPlayerFound = window.QPoker.OnPlayerFound || [];
window.QPoker.InitLogin = onLoginClick;
