import * as React from "react";

import { classNames } from "../utils";
import { Game } from "../objects/Game";

type SideBarProps = {
    game?: Game;
    playerID: string;
    playerToken: string;
    sendAction: (action: any) => void;
}

type SideBarState = {
    selectedTab: string;
}

export class SideBar extends React.Component<SideBarProps, SideBarState> {
    constructor(props: any) {
        super(props);
        this.state = { selectedTab: "history" };
    }

    public getNavOptions(): string[] {
        let options = ['history', 'chat'];

        if (this.props.game.owner_id.toString() == this.props.playerID.toString()) {
            options.push('admin');
        }

        return options
    }

    public navSelect(event: any) {
        this.setState({ selectedTab: event.target.innerHTML });

        // TODO: handle newly selected tab
        event.stopPropagation();
    }

    public render() {
        return <div className="sidebar">
            <div className="row sidebar-nav">
                {this.getNavOptions().map((label) => {
                    return <a href="#"
                                className={classNames("btn-flat", { underline: this.state.selectedTab === label })}
                                onClick={this.navSelect.bind(this)} >
                        {label}
                    </a>;
                })}
            </div>
        </div>;
    }
}