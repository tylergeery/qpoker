import * as React from "react";
import * as ReactDOM from "react-dom";

import { Nav } from "./components/Nav";
import { LoginModal } from "./components/modals/Login";
import { RegistrationModal } from "./components/modals/Registration";
import { LoginRequest } from "./requests/login";
import { RegistrationRequest } from "./requests/registration";
import { userStorage } from "./utils/storage";

var registrationModal: RegistrationModal;
var loginModal: LoginModal;
var nav: Nav;

async function submitLogin (loginData: object) {
    const req = new LoginRequest();
    const player = await req.request(loginData)

    if (!player || !player.id) {
        debugger; // TODO: handle errors
        return
    }

    userStorage.setUser(player);
    nav.setState({ player });
    loginModal.closeModal();
}

async function submitRegistration (loginData: object) {
    const req = new RegistrationRequest();
    const player = await req.request(loginData)

    if (!player || !player.id) {
        debugger; // TODO: handle errors
        return
    }

    userStorage.setUser(player);
    nav.setState({ player });
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
