import * as React from "react";

import { Player } from "../objects/Player";
import { userStorage } from "../utils/storage";
import { GetPlayerRequest } from "../requests/getPlayer";

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
        const userID = userStorage.getID();
        const userToken = userStorage.getToken();

        if (userID && userToken) {
            // Send get user request
            let req = new GetPlayerRequest<Player>();
            let player = await req.request({id: userID, userToken});
            if (req.success) {
                this.setState({player});
            }
        }
    }

    public logout() {
        userStorage.removePlayer(this.state.player);
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