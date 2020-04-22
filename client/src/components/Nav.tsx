import * as React from "react";

import { Player } from "../objects/Player";
import { getPlayer } from "../helpers/player";
import { userStorage } from "../utils/storage";

type NavProps = {
    onLoginClicked: (event: any) => void;
    onRegisterClicked: (event: any) => void;
}

interface NavState {
    player?: Player;
}

// State is never set so we use the '{}' type.
export class Nav extends React.Component<NavProps, NavState> {
    constructor(props: any) {
        super(props)
        this.state = {player: null};
    }
    public async componentDidMount() {
        let player = await getPlayer();

        if (player) {
            this.setState({player});
        }
    }

    public logout() {
        userStorage.removePlayer();
        this.setState({player: null});
    }

    public render() {
        return this.state && this.state.player ? (
            <li><a onClick={this.logout.bind(this)} href="#logout">Logout</a></li>
        ) : (
            <span>
                <li><a onClick={this.props.onLoginClicked} href="#login">Log In</a></li>
                <li><a onClick={this.props.onRegisterClicked}href="#signup">Sign Up</a></li>
            </span>
        );
    }
}